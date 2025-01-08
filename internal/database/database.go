package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// Connect establishes a connection to the given Postgres database and returns a connection pool to be used for further access.
func Connect(ctx context.Context, cfg *Config, debugMode bool) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode,
	)

	pgxPoolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse pgx pool config")
	}
	pgxPoolCfg.ConnConfig.Tracer = newTracer(debugMode)
	pgxPoolCfg.MaxConns = cfg.MaxConns
	pgxPoolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgxPoolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	connPool, err := pgxpool.NewWithConfig(ctx, pgxPoolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create db connection pool")
	}

	// The database is pinged to ensure the connection was established.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := connPool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "error pinging postgres")
	}

	return connPool, nil
}
