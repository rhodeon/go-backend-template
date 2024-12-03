package test_utils

import (
	"path"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rhodeon/go-backend-template/internal/database"
)

type Config struct {
	PostgresContainer string `env:"TEST_POSTGRES_CONTAINER" envDefault:"required"`

	Database *database.Config
}

var config *Config

func parseConfig() (*Config, error) {
	if err := godotenv.Load(path.Join(projectRootDir, ".env")); err != nil {
		return nil, errors.Wrap(err, "failed to load .env file")
	}

	cfg := env.Must(env.ParseAs[Config]())

	cfg.Database = &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "airgateway",
		Pass:     "password",
		Name:     "agw_test",
		SslMode:  "disable",
		MaxConns: 1,
	}

	return &cfg, nil
}
