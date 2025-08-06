package postgres

import (
	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
)

type Repository struct {
	Users *dbusers.Queries
}

func NewRepository() *Repository {
	return &Repository{
		Users: dbusers.New(),
	}
}
