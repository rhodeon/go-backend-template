package main

import (
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/repositories"
	"sync"
)

func main() {
	cfg := internal.ParseConfig()
	logger := log.NewLogger(cfg.DebugMode)
	repos := repositories.New()

	dbConfig := database.Config(cfg.Database)
	dbPool, err := database.Connect(&dbConfig, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}

	app := internal.NewApplication(cfg, logger, repos, dbPool)

	// A waitgroup is established to ensure background tasks are completed before shutting down the server.
	backgroundWg := &sync.WaitGroup{}

	// Start server.
	err = serveApi(app, backgroundWg)
	if err != nil {
		panic(err)
	}
}
