package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rhodeon/go-backend-template/internal/config"
	"log/slog"
	"time"
)

// Connect establishes a connection to the given Postgres database and returns a connection pool to be used for further access.
func Connect(cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	dbCfg := cfg.Database

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SslMode,
	)

	pgxPoolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse pgx pool config")
	}
	pgxPoolCfg.ConnConfig.Tracer = newTracer(logger, cfg.DebugMode)
	pgxPoolCfg.MaxConns = dbCfg.MaxConns
	pgxPoolCfg.MaxConnLifetime = dbCfg.MaxConnLifetime
	pgxPoolCfg.MaxConnIdleTime = dbCfg.MaxConnIdleTime

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create db connection pool")
	}

	// The database is pinged to ensure the connection was established.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := connPool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "error pinging postgres")
	}

	return connPool, nil
}
