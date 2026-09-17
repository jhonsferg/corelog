package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/user/application/port"
	"github.com/jhonsferg/corelog/internal/user/domain"
)

const selectColumns = `id, name, email, password_hash, role, team_id, created_at, updated_at`

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ port.Repository = (*UserRepository)(nil)

func (r *UserRepository) Save(ctx context.Context, u *domain.User) error {
	row := fromDomain(u)

	const query = `
		INSERT INTO users (id, name, email, password_hash, role, team_id, created_at, updated_at)
		VALUES (:id, :name, :email, :password_hash, :role, :team_id, :created_at, :updated_at)`

	if _, err := r.db.NamedExecContext(ctx, query, row); err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to save user", err)
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var row userRow

	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1`, selectColumns)

	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, apperror.Wrap(apperror.KindInternal, "failed to find user", err)
	}

	return row.toDomain(), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var row userRow

	query := fmt.Sprintf(`SELECT %s FROM users WHERE email = $1`, selectColumns)

	if err := r.db.GetContext(ctx, &row, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, apperror.Wrap(apperror.KindInternal, "failed to find user", err)
	}

	return row.toDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	row := fromDomain(u)

	const query = `
		UPDATE users
		SET name = :name, email = :email, password_hash = :password_hash, team_id = :team_id, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to update user", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to update user", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.db.GetContext(ctx, &count, `SELECT COUNT(1) FROM users`); err != nil {
		return 0, apperror.Wrap(apperror.KindInternal, "failed to count users", err)
	}
	return count, nil
}

func (r *UserRepository) Search(ctx context.Context, filter port.SearchFilter) ([]*domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE 1 = 1`, selectColumns)
	args := []any{}

	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", len(args), len(args))
	}
	if filter.TeamID != nil {
		args = append(args, *filter.TeamID)
		query += fmt.Sprintf(" AND team_id = $%d", len(args))
	}

	query += " ORDER BY name ASC LIMIT 20"

	var rows []userRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, apperror.Wrap(apperror.KindInternal, "failed to search users", err)
	}

	users := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, row.toDomain())
	}
	return users, nil
}
