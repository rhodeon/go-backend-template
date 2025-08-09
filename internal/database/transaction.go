package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Tx represents a wrapper around a pgx.Tx, providing methods to interact with database transactions.
// Importantly, it doesn't expose methods to commit or rollback. Those are delegated to be done at the same level where the transaction is created via Db.BeginTx.
// This is done to eliminate the possibility of terminating a transaction down the call stack from where it was created and attempting to use it afterwards.
type Tx struct {
	innerTx pgx.Tx
}

// Savepoint creates an anonymous savepoint and returns its rollback function.
// Unlike Db.BeginTx, the rollback returned here should never be deferred as that will always roll back the parent transaction.
// Instead, it should be used in the event that a recoverable error occurs during a transaction.
// Savepoints are not cheap at a large scale, so this should be used sparingly.
// Where possible, prefer handling recoverable errors with ON CONFLICT instead.
// More context on their cost: https://postgres.ai/blog/20210831-postgresql-subtransactions-considered-harmful
func (tx *Tx) Savepoint(ctx context.Context) (func(ctx context.Context) error, error) {
	sp, err := tx.innerTx.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting transaction savepoint: %w", err)
	}
	return tx.rollbackSavepoint(sp), nil
}

func (tx *Tx) rollbackSavepoint(sp pgx.Tx) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// If the transaction is already closed, the error can be ignored.
		if err := sp.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return fmt.Errorf("rolling back transaction savepoint: %w", err)
		}
		return nil
	}
}

// Exec fulfills the DBTX interface needed for sqlc `:exec` operations.
func (tx *Tx) Exec(ctx context.Context, s string, a ...any) (pgconn.CommandTag, error) {
	return tx.innerTx.Exec(ctx, s, a...) //nolint:wrapcheck
}

// Query fulfills the DBTX interface needed for sqlc `:many` operations.
func (tx *Tx) Query(ctx context.Context, s string, a ...any) (pgx.Rows, error) {
	return tx.innerTx.Query(ctx, s, a...) //nolint:wrapcheck
}

// QueryRow fulfills the DBTX interface needed for sqlc `:one` operations.
func (tx *Tx) QueryRow(ctx context.Context, s string, a ...any) pgx.Row {
	return tx.innerTx.QueryRow(ctx, s, a...) //nolint:wrapcheck
}

// CopyFrom fulfills the DBTX interface needed for sqlc `:copyfrom` operations.
func (tx *Tx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return tx.innerTx.CopyFrom(ctx, tableName, columnNames, rowSrc) //nolint:wrapcheck
}
