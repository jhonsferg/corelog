package database

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"

	"github.com/jhonsferg/corelog/internal/platform/config"
	sqlitemigrations "github.com/jhonsferg/corelog/migrations/sqlite"
)

func init() {
	sqlx.BindDriver("sqlite", sqlx.QUESTION)
}

func NewSQLite(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	path := cfg.SQLitePath()

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("database: create sqlite directory: %w", err)
		}
	}

	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", path)

	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("database: connect: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	if err := applySQLiteMigrations(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func applySQLiteMigrations(ctx context.Context, db *sqlx.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)`); err != nil {
		return fmt.Errorf("database: create schema_migrations table: %w", err)
	}

	entries, err := fs.ReadDir(sqlitemigrations.FS, ".")
	if err != nil {
		return fmt.Errorf("database: read sqlite migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		version := strings.SplitN(name, "_", 2)[0]

		var applied int
		if err := db.GetContext(ctx, &applied, `SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, version); err != nil {
			return fmt.Errorf("database: check sqlite migration %s: %w", name, err)
		}
		if applied > 0 {
			continue
		}

		content, err := fs.ReadFile(sqlitemigrations.FS, name)
		if err != nil {
			return fmt.Errorf("database: read sqlite migration %s: %w", name, err)
		}

		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("database: begin sqlite migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("database: apply sqlite migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("database: record sqlite migration %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("database: commit sqlite migration %s: %w", name, err)
		}
	}

	return nil
}
