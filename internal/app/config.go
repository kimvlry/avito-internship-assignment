package app

import (
    "fmt"
    "github.com/go-playground/validator/v10"
    "github.com/ilyakaznacheev/cleanenv"
    "github.com/joho/godotenv"
    "time"
)

type Config struct {
    AppMode string `env:"APP_MODE" env-default:"dev"`

    Postgres PostgresConfig
    Http     HttpConfig
}

type PostgresConfig struct {
    Host     string `env:"PG_HOST" env-default:"localhost" validate:"required"`
    Port     string `env:"PG_PORT" env-default:"5432" validate:"required"`
    User     string `env:"PG_USER" env-default:"user" validate:"required"`
    Password string `env:"PG_PASSWORD" env-default:"pass" validate:"required"`
    DBName   string `env:"PG_DB" env-default:"avito-fall-2025" validate:"required"`
    SSLMode  string `env:"PG_SSLMODE" env-default:"disable" validate:"oneof=disable allow prefer require verify-ca verify-full"`
}

type HttpConfig struct {
    Port         string        `env:"HTTP_PORT" validate:"required"`
    ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s" validate:"required"`
    WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s" validate:"required"`
    IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"5s" validate:"required"`
    JwtSecret    string        `env:"JWT_SECRET" validate:"required"`
}

func (h HttpConfig) Addr() string {
    return ":" + h.Port
}

func LoadConfig() (*Config, error) {
    _ = godotenv.Load()

    var cfg Config
    if err := cleanenv.ReadEnv(&cfg); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }

    validate := validator.New()
    if err := validate.Struct(cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    return &cfg, nil
}

func (p PostgresConfig) GetConnString() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode,
    )
}
