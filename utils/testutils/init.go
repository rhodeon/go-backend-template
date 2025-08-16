package testutils

import (
	"context"
	"log"

	"github.com/rhodeon/go-backend-template/repositories/cache/redis"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
	"github.com/testcontainers/testcontainers-go"
)

type ContainerOpts struct {
	Postgres bool
	Redis    bool
}

func SetupContainers(ctx context.Context, opts ContainerOpts) (func(context.Context) error, error) {
	var err error
	var containers []testcontainers.Container

	if config, err = parseConfig(); err != nil {
		log.Fatal(err)
	}

	if opts.Postgres {
		postgresContainer, err := postgres.SetupTestContainer(ctx, config.PostgresImage)
		if err != nil {
			log.Fatal(err)
		}
		containers = append(containers, postgresContainer)
	}

	if opts.Redis {
		redisContainer, err := redis.SetupTestContainer(ctx, config.RedisImage)
		if err != nil {
			log.Fatal(err)
		}
		containers = append(containers, redisContainer)
	}

	return terminateContainers(containers), nil
}

func terminateContainers(containers []testcontainers.Container) func(context.Context) error {
	return func(ctx context.Context) error {
		for _, container := range containers {
			if err := container.Terminate(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}
