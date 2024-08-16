package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rhodeon/go-backend-template/repositories/database/implementation/posts"
	"github.com/rhodeon/go-backend-template/repositories/database/implementation/users"
)

type Repository struct {
	Users users.Querier
	Posts posts.Querier
}

func NewRepository() *Repository {
	return &Repository{
		Users: users.New(),
		Posts: posts.New(),
	}
}

type Transaction interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}
