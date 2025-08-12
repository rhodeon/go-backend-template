package services

import (
	"github.com/rhodeon/go-backend-template/repositories"
)

type Services struct {
	Auth *Auth
	User *User
}

func New(repos *repositories.Repositories, cfg *Config) *Services {
	return &Services{
		newAuth(repos, cfg),
		newUser(repos, cfg),
	}
}

// service is used as a branded type to have uniform properties and resources across the various services.
type service struct {
	repos *repositories.Repositories
	cfg   *Config
}

func newService(repos *repositories.Repositories, cfg *Config) *service {
	return &service{
		repos,
		cfg,
	}
}
