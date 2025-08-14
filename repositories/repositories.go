package repositories

import (
	"github.com/rhodeon/go-backend-template/repositories/cache"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
	"github.com/rhodeon/go-backend-template/repositories/email"
)

type Repositories struct {
	Database *postgres.Database
	Cache    cache.Cache
	Email    email.Email
}

func New(cacheRepo cache.Cache, emailRepo email.Email) *Repositories {
	return &Repositories{
		postgres.New(),
		cacheRepo,
		emailRepo,
	}
}
