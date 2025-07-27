package database

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/helpers"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// BeginTransaction starts and returns a new transaction from the given pool along with its associated commit and rollback resolver functions.
// Both resolvers are returned as guardrails to reduce the chance of forgetting to commit/rollback after the transaction is done.
// The rollback is safe to call after commits and should always be deferred after calling BeginTransaction as a fail-safe.
func BeginTransaction(ctx context.Context, dbPool *pgxpool.Pool) (pgx.Tx, commitResolver, rollbackResolver, error) {
	dbTx, err := dbPool.Begin(ctx)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "starting database transaction")
	}

	return dbTx, commitTransaction(dbTx), rollbackTransaction(dbTx), nil
}

type (
	commitResolver   func(ctx context.Context) error
	rollbackResolver func(ctx context.Context)
)

func commitTransaction(tx pgx.Tx) commitResolver {
	return func(ctx context.Context) error {
		// If the transaction is already closed, the error can be ignored.
		if err := tx.Commit(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return errors.Wrap(err, "commiting database transaction")
		}
		return nil
	}
}

func rollbackTransaction(tx pgx.Tx) rollbackResolver {
	return func(ctx context.Context) {
		// If the transaction is already closed, the error can be ignored.
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			helpers.ContextGetLogger(ctx).Error("Rolling back database transaction", slog.Any(log.AttrError, err))
		}
	}
}
