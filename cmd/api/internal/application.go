package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rhodeon/go-backend-template/domain/services"
	"log/slog"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config   *Config
	Logger   *slog.Logger
	DbPool   *pgxpool.Pool
	Services *services.Services
}

func NewApplication(cfg *Config, logger *slog.Logger, dbPool *pgxpool.Pool, services *services.Services) *Application {
	return &Application{
		Config:   cfg,
		Logger:   logger,
		DbPool:   dbPool,
		Services: services,
	}
}
