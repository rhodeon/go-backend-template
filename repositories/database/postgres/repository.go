package postgres

import (
	"context"

	pgposts "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/posts"
	pgusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	Users *pgusers.Queries
	Posts *pgposts.Queries
}

func NewRepository() *Repository {
	return &Repository{
		Users: pgusers.New(),
		Posts: pgposts.New(),
	}
}

type Transaction interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
