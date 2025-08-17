package redis

import (
	"context"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/go-errors/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rhodeon/go-backend-template/repositories/cache"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// containerName is the name of the global Redis container to be shared across all test packages.
// The UUID suffix is meant to reduce (practically eliminate) the chances of collision with another container.
const containerName = "gobt-redis-8b44d6d1"

// maxDatabases defines the maximum number of Redis databases.
// It also limits the number of concurrent tests which depend on Redis and should be increased accordingly.
// For the best effect, it should be set to at least 10x the number of concurrent tests to reduce the chances of repetition
// when generating the database number per-test.
// 1000 was arbitrarily chosen as a reasonable value.
const maxDatabases = 1000

var testConfig = &Config{
	Host: "localhost",
}

// SetupTestContainer establishes a global Redis instance in a container to be used for testing.
// This container is truly "global". In order words, a single container is shared/reused across all test packages in the codebase.
func SetupTestContainer(ctx context.Context, image string, projectRootDir string) error {
	redisContainer, err := tcredis.Run(ctx,
		image,
		testcontainers.WithReuseByName(containerName),
		tcredis.WithConfigFile(filepath.Join(projectRootDir, "testdata", "redis.conf")),
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

var errReservedDatabase = errors.New("redis database is reserved by another test")

// NewTestCache creates and returns a cache connection to a container established by SetupTestContainer.
// Each test connects to a different database to allow independent operations.
func NewTestCache(ctx context.Context) (cache.Cache, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for attempt := 0; attempt < maxDatabases; attempt++ {
		dbNumber := rand.Intn(maxDatabases)
		redisConfig := *testConfig
		redisConfig.Database = dbNumber

		// The first test to connect to a database gets an exclusive lock on it.
		redisConfig.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
			if !cn.SetNX(ctx, "reserved_for_test", 1, 0).Val() {
				return errors.New(errReservedDatabase)
			}
			return nil
		}

		redisCache, err := New(ctx, &redisConfig)
		if err != nil {
			switch {
			case errors.Is(err, errReservedDatabase):
				// The database is reserved, so another one is tried.
				continue

			default:
				return nil, errors.Errorf("connecting redis: %w", err)
			}

		}

		return redisCache, nil
	}

	return nil, errors.Errorf("failed to find an available redis database after %d tries", maxDatabases)
}
