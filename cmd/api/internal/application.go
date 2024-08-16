package internal

import (
	"log/slog"
)

// Application houses common resources used at various points in the API as a means of dependency injection.
type Application struct {
	Config *Config
	Logger *slog.Logger
}

func NewApplication(cfg *Config, logger *slog.Logger) *Application {
	app := &Application{
		Config: cfg,
		Logger: logger,
	}

	return app
}
