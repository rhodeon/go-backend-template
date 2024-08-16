package services

import (
	"github.com/rhodeon/go-backend-template/repositories"
	"log/slog"
)

type Services struct {
	User *User
}

func New(repos *repositories.Repositories, logger *slog.Logger) *Services {
	return &Services{
		User: newUser(repos, logger),
	}
}

// service is used as a branded type to have uniform properties and resources across the various services.
type service struct {
	repos  *repositories.Repositories
	logger *slog.Logger
}

func newService(repos *repositories.Repositories, logger *slog.Logger) *service {
	return &service{
		repos:  repos,
		logger: logger,
	}
}
