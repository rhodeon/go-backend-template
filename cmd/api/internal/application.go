package internal

import (
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/services"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config   *Config
	Logger   *slog.Logger
	Db       *database.Db
	Services *services.Services
}

func NewApplication(cfg *Config, logger *slog.Logger, db *database.Db, svcs *services.Services) *Application {
	return &Application{
		Config:   cfg,
		Logger:   logger,
		Db:       db,
		Services: svcs,
	}
}
