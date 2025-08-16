package testutils

import (
	"path"

	"github.com/rhodeon/go-backend-template/utils/testutils/temp"

	"github.com/caarlos0/env/v11"
	"github.com/go-errors/errors"
	"github.com/joho/godotenv"
)

type Config struct {
	PostgresImage string `env:"TEST_POSTGRES_IMAGE,required"`
	RedisImage    string `env:"TEST_REDIS_IMAGE,required"`
}

var config *Config

func parseConfig() (*Config, error) {
	if err := godotenv.Load(path.Join(temp.ProjectRootDir, ".env")); err != nil {
		return nil, errors.Errorf("loading .env file: %w", err)
	}

	cfg := env.Must(env.ParseAs[Config]())

	return &cfg, nil
}
