package redis

import (
	"context"
	"time"

	"github.com/rhodeon/go-backend-template/repositories/cache"

	"github.com/go-errors/errors"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// containerName is the name of the global Redis container to be shared across all test packages.
// The UUID suffix is meant to reduce (practically eliminate) the chances of collision with another container.
const containerName = "gobt-redis-8b44d6d1"

var testConfig *Config

// SetupTestContainer establishes a global Redis instance in a container to be used for testing.
// This container is truly "global". In order words, a single container is shared/reused across all test packages in the codebase.
func SetupTestContainer(ctx context.Context, image string) error {
	redisContainer, err := tcredis.Run(ctx,
		image,
		testcontainers.WithReuseByName(containerName),
	)
	if err != nil {
		return errors.Errorf("creating Redis container instance: %w", err)
	}

	mappedPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return errors.Errorf("getting mapped Redis container ports: %w", err)
	}

	testConfig = &Config{
		Host: "localhost",
		Port: mappedPort.Int(),
	}

	if err = redisContainer.Start(ctx); err != nil {
		return errors.Errorf("starting Redis container: %w", err)
	}

	return nil
}

// NewTestCache creates and returns a cache connection to a container established by SetupTestContainer.
func NewTestCache(ctx context.Context) (cache.Cache, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	redisCache, err := New(ctx, testConfig)
	if err != nil {
		return nil, errors.Errorf("connecting redis: %w", err)
	}

	return redisCache, nil
}
