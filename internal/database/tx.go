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

// BeginTransaction starts a transaction and returns the new transaction along with the rollback function to be deferred.
func BeginTransaction(ctx context.Context, dbPool *pgxpool.Pool) (pgx.Tx, func(ctx context.Context, tx pgx.Tx), error) {
	dbTx, err := dbPool.Begin(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "starting database transaction")
	}

	return dbTx, ResolveTransactionWithRollback, nil
}

// ResolveTransactionWithRollback acts as a failsafe, preventing any missed rollbacks from holding on to a transaction.
// It is safe to call after commits and should always be deferred when creating a transaction
func ResolveTransactionWithRollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		helpers.GetContextLogger(ctx).ErrorContext(ctx, "Error resolving database transaction connection", slog.Any(log.AttrError, err))
	}
}
