package testutils

import (
	"context"

	"github.com/rhodeon/go-backend-template/repositories/cache/redis"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"

	"github.com/go-errors/errors"
)

// ContainerOpts defines options for setting up containerised resources for integration tests.
// The state of the different fields indicates if their corresponding containers should be created.
type ContainerOpts struct {
	Postgres bool
	Redis    bool
}

// SetupContainers sets up the resources needed before running integration tests.
// Since the containers are shared between all test packages, clean-up cannot be easily orchestrated internally.
// Instead, the Testcontainers reaper process automatically handles the clean-up after all tests are done.
// Reaper reference: https://golang.testcontainers.org/features/garbage_collector/#ryuk
func SetupContainers(ctx context.Context, opts ContainerOpts) error {
	if opts.Postgres {
		if err := postgres.SetupTestContainer(ctx, config.PostgresImage, projectRootDir); err != nil {
			return errors.Errorf("setting up postgres container: %w", err)
		}
	}

	if opts.Redis {
		if err := redis.SetupTestContainer(ctx, config.RedisImage, projectRootDir); err != nil {
			return errors.Errorf("setting up redis container: %w", err)
		}
	}

	return nil
}
