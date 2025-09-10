package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/utils/contextutils"

	"github.com/exaring/otelpgx"
	"github.com/go-errors/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Db abstracts the underlying database pool and provides helper methods for conveniently creating and managing transactions.
// Another reason for the abstraction is to encourage the consistent use of transactions for all database operations.
// Having this pattern makes it easier to extend flows (like some which where originally read-only).
// Postgres always creates an implicit transaction in any case, so making it explicit only has a negligible cost of an extra round-trip.
// While the Db.Pool method exposes the underlying pool, that should be reserved only for testing and a few exceptional cases.
type Db struct {
	pool *pgxpool.Pool
}

// Connect establishes a connection to the given Postgres database and returns a connection pool to be used for further access.
func Connect(ctx context.Context, cfg *Config, debugMode bool) (*Db, func(), error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode,
	)

	pgxPoolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, errors.Errorf("parsing pgx pool config: %w", err)
	}

	tracers := multitracer.New(
		newTracer(debugMode),
		otelpgx.NewTracer(),
	)
	pgxPoolCfg.ConnConfig.Tracer = tracers

	pgxPoolCfg.MaxConns = cfg.MaxConns
	pgxPoolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgxPoolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	connPool, err := pgxpool.NewWithConfig(ctx, pgxPoolCfg)
	if err != nil {
		return nil, nil, errors.Errorf("creating db connection pool: %w", err)
	}

	// The database is pinged to ensure the connection is established.
	maxAttempts := 10
	delay := 2 * time.Second

	var errPing error
	for i := 0; i < maxAttempts; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		errPing = connPool.Ping(pingCtx)
		cancel()

		if errPing == nil {
			return &Db{connPool}, connPool.Close, nil
		}

		time.Sleep(delay)
	}

	return nil, nil, errors.Errorf("pinging postgres after %d attempts: %w", maxAttempts, errPing)
}

// TxOptions wraps the pgx-specific transaction options and leaves room for extending with other options.
// For example, a flag can be added to prevent updating audit logs for certain operations.
type TxOptions struct {
	pgx.TxOptions
}

// BeginTx starts and returns a new transaction from the given pool along with its associated commit and rollback resolver functions.
// Both resolvers are returned as guardrails to reduce the chance of forgetting to commit/rollback after the transaction is done.
// The rollback is safe to call after commits and should always be deferred after calling BeginTx as a fail-safe.
func (p *Db) BeginTx(ctx context.Context, opts ...TxOptions) (*Tx, commitResolver, rollbackResolver, error) {
	var txOptions pgx.TxOptions
	if len(opts) > 0 {
		txOptions = opts[0].TxOptions
	}

	dbTx, err := p.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, nil, nil, errors.Errorf("starting database transaction: %w", err)
	}

	tx := &Tx{dbTx}
	return tx, p.commitTransaction(tx), p.rollbackTransaction(tx), nil
}

type (
	commitResolver   func(ctx context.Context) error
	rollbackResolver func(ctx context.Context)
)

func (p *Db) commitTransaction(tx *Tx) commitResolver {
	return func(ctx context.Context) error {
		// If the transaction is already closed, the error can be ignored.
		if err := tx.innerTx.Commit(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return errors.Errorf("commiting database transaction: %w", err)
		}
		return nil
	}
}

func (p *Db) rollbackTransaction(tx *Tx) rollbackResolver {
	return func(ctx context.Context) {
		// If the transaction is already closed, the error can be ignored.
		if err := tx.innerTx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			contextutils.GetLogger(ctx).Error("Rolling back database transaction", slog.Any(log.AttrError, err))
		}
	}
}

// Pool exposes the underlying connection pool. This is useful for tests where vetting the results of database
// operations can be done directly without introducing the boilerplate of managing transactions.
func (p *Db) Pool() *pgxpool.Pool {
	return p.pool
}
