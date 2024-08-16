package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rhodeon/go-backend-template/repositories"
	"log/slog"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config     *Config
	Logger     *slog.Logger
	Repository *repositories.Repositories
	DbPool     *pgxpool.Pool
}

func NewApplication(cfg *Config, logger *slog.Logger, repo *repositories.Repositories, dbPool *pgxpool.Pool) *Application {
	return &Application{
		Config:     cfg,
		Logger:     logger,
		Repository: repo,
		DbPool:     dbPool,
	}
}
