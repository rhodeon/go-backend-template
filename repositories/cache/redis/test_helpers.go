package redis

import (
	"context"
	"time"

	"github.com/rhodeon/go-backend-template/repositories/cache"

	"github.com/go-errors/errors"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var testConfig *Config

// SetupTestContainer establishes a Redis instance in a container to be used for testing.
func SetupTestContainer(ctx context.Context, image string) (*tcredis.RedisContainer, error) {
	redisContainer, err := tcredis.Run(ctx,
		image,
	)
	if err != nil {
		return nil, errors.Errorf("creating Redis container instance: %w", err)
	}

	mappedPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return nil, errors.Errorf("getting mapped Redis container ports: %w", err)
	}

	testConfig = &Config{
		Host: "localhost",
		Port: mappedPort.Int(),
	}

	if err = redisContainer.Start(ctx); err != nil {
		return nil, errors.Errorf("starting Redis container: %w", err)
	}

	return redisContainer, nil
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
