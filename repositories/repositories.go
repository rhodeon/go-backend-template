package repositories

import (
	"github.com/rhodeon/go-backend-template/repositories/database"
)

type Repositories struct {
	Database *database.Repository
}

func New() *Repositories {
	return &Repositories{
		Database: database.NewRepository(),
	}
}
