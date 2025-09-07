package postgres

import (
	dbpetcategories "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/petcategories"
	dbpets "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/pets"
	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
)

type Database struct {
	Users         *dbusers.Queries
	Pets          *dbpets.Queries
	PetCategories *dbpetcategories.Queries
}

func New() *Database {
	return &Database{
		dbusers.New(),
		dbpets.New(),
		dbpetcategories.New(),
	}
}
