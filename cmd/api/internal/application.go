package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rhodeon/go-backend-template/domain/services"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config   *Config
	DbPool   *pgxpool.Pool
	Services *services.Services
}

func NewApplication(cfg *Config, dbPool *pgxpool.Pool, services *services.Services) *Application {
	return &Application{
		Config:   cfg,
		DbPool:   dbPool,
		Services: services,
	}
}
