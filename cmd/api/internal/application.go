package internal

import (
	"github.com/rhodeon/go-backend-template/internal/config"
	"log/slog"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config config.Config
	Logger *slog.Logger
}

func NewApplication(cfg *config.Config, logger *slog.Logger) *Application {
	app := &Application{
		Config: *cfg,
		Logger: logger,
	}

	return app
}
