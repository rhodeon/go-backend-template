package postgres

import (
	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
)

type Database struct {
	Users *dbusers.Queries
}

func NewDatabase() *Database {
	return &Database{
		Users: dbusers.New(),
	}
}
