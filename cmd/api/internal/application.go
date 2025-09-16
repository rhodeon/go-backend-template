package internal

import (
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/services"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config   *Config
	Db       *database.Db
	Services *services.Services
}

func NewApplication(cfg *Config, db *database.Db, svcs *services.Services) *Application {
	return &Application{
		Config:   cfg,
		Db:       db,
		Services: svcs,
	}
}
