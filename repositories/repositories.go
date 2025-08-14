package repositories

import (
	"github.com/rhodeon/go-backend-template/repositories/cache"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
)

type Repositories struct {
	Database *postgres.Database
	Cache    cache.Cache
}

func New(cacheRepo cache.Cache) *Repositories {
	return &Repositories{
		postgres.NewDatabase(),
		cacheRepo,
	}
}
