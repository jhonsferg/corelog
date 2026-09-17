package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Features FeatureConfig
}

type ServerConfig struct {
	Port            string        `env:"SERVER_PORT" envDefault:"8080"`
	Env             string        `env:"APP_ENV" envDefault:"development"`
	AllowedOrigins  []string      `env:"CORS_ALLOWED_ORIGIN" envDefault:"http://localhost:5173" envSeparator:","`
	ReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"15s"`
}

type DatabaseConfig struct {
	Driver          string        `env:"DB_DRIVER" envDefault:"postgres"`
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            string        `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"corelog"`
	Password        string        `env:"DB_PASSWORD" envDefault:"corelog"`
	Name            string        `env:"DB_NAME" envDefault:"corelog"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func (d DatabaseConfig) SQLitePath() string {
	return filepath.Join("data", d.Name+".db")
}

type AuthConfig struct {
	JWTSecret     string        `env:"JWT_SECRET" envDefault:"change-me-in-production"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" envDefault:"24h"`
}

type FeatureConfig struct {
	TicketRepository string `env:"TICKET_REPOSITORY" envDefault:"postgres"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}

	return cfg, nil
}
