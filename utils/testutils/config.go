package testutils

import (
	"fmt"
	"path"

	"github.com/rhodeon/go-backend-template/internal/database"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	PostgresContainer string `env:"TEST_POSTGRES_CONTAINER,required"`

	Database *database.Config
}

var config *Config

func parseConfig() (*Config, error) {
	if err := godotenv.Load(path.Join(projectRootDir, ".env")); err != nil {
		return nil, fmt.Errorf("loading .env file: %w", err)
	}

	cfg := env.Must(env.ParseAs[Config]())

	cfg.Database = &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "test_user",
		Pass:     "test_pass",
		Name:     "test_db",
		SslMode:  "disable",
		MaxConns: 1,
	}

	return &cfg, nil
}
