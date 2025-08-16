package testutils

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/rhodeon/go-backend-template/repositories/cache/redis"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
	"github.com/testcontainers/testcontainers-go"
)

// ContainerOpts defines options for setting up containerised resources for integration tests.
// The state of the different fields indicates if their corresponding containers should be created.
type ContainerOpts struct {
	Postgres bool
	Redis    bool
}

// SetupContainers sets up the resources needed before running integration tests.
func SetupContainers(ctx context.Context, opts ContainerOpts) (func(context.Context) error, error) {
	var containers []testcontainers.Container

	if opts.Postgres {
		postgresContainer, err := postgres.SetupTestContainer(ctx, config.PostgresImage, projectRootDir)
		if err != nil {
			return nil, errors.Errorf("setting up postgres container: %w", err)
		}
		containers = append(containers, postgresContainer)
	}

	if opts.Redis {
		redisContainer, err := redis.SetupTestContainer(ctx, config.RedisImage)
		if err != nil {
			return nil, errors.Errorf("setting up redis container: %w", err)
		}
		containers = append(containers, redisContainer)
	}

	return terminateContainers(containers), nil
}

func terminateContainers(containers []testcontainers.Container) func(context.Context) error {
	return func(ctx context.Context) error {
		for _, container := range containers {
			if err := container.Terminate(ctx); err != nil {
				return errors.Errorf("terminating container: %w", err)
			}
		}

		return nil
	}
}
