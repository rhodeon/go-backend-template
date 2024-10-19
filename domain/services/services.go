package services

import (
	"github.com/rhodeon/go-backend-template/repositories"
)

type Services struct {
	User *User
}

func New(repos *repositories.Repositories) *Services {
	return &Services{
		User: newUser(repos),
	}
}

// service is used as a branded type to have uniform properties and resources across the various services.
type service struct {
	repos *repositories.Repositories
}

func newService(repos *repositories.Repositories) *service {
	return &service{
		repos: repos,
	}
}
