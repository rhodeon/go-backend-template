package redis

import (
	"context"
	"log/slog"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/repositories/cache"

	"github.com/go-errors/errors"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	imageName = "redis:8.2.0"

	// containerName is the name of the global Redis container to be shared across all test packages.
	// The UUID suffix is meant to reduce (practically eliminate) the chances of collision with another container.
	containerName = "gobt-redis-8b44d6d1"

	// maxDatabases defines the maximum number of Redis databases.
	// It also limits the number of concurrent tests which depend on Redis and should be increased accordingly.
	// For the best effect, it should be set to at least 10x the number of concurrent tests to reduce the chances of repetition
	// when generating the database number per-test.
	// 1000 was arbitrarily chosen as a reasonable value.
	maxDatabases = 1000
)

var testConfig = &Config{
	Host: "localhost",
}

// SetupTestContainer establishes a global Redis instance in a container to be used for testing.
// This container is truly "global". In other words, a single container is shared/reused across all test packages in the codebase.
func SetupTestContainer(ctx context.Context, projectRootDir string) error {
	redisContainer, err := tcredis.Run(ctx,
		imageName,
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

	testConfig.Port = mappedPort.Int()
	slog.InfoContext(ctx, "Redis test container is ready", slog.Int(log.AttrPort, testConfig.Port))

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

		slog.InfoContext(
			ctx,
			"Connected to Redis test cache",
			slog.Int(log.AttrDatabase, dbNumber),
			slog.Int(log.AttrPort, testConfig.Port),
		)

		return redisCache, nil
	}

	return nil, errors.Errorf("failed to find an available redis database after %d tries", maxDatabases)
}
