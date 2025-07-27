package repositories

import (
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
)

type Repositories struct {
	Database *postgres.Repository
}

func New() *Repositories {
	return &Repositories{
		Database: postgres.NewRepository(),
	}
}
