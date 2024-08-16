package main

import (
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/services"
	"sync"
)

func main() {
	cfg := internal.ParseConfig()
	logger := log.NewLogger(cfg.DebugMode)
	logger.Debug("HEYYYYYYY")

	dbConfig := database.Config(cfg.Database)
	dbPool, err := database.Connect(&dbConfig, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}

	repos := repositories.New()
	services := services.New(repos, logger)

	app := internal.NewApplication(cfg, logger, dbPool, services)

	// A waitgroup is established to ensure background tasks are completed before shutting down the server.
	backgroundWg := &sync.WaitGroup{}

	// Start server.
	err = serveApi(app, backgroundWg)
	if err != nil {
		panic(err)
	}
}
