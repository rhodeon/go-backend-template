package main

import (
	"sync"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/server"
	"github.com/rhodeon/go-backend-template/domain/services"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/repositories"
)

func main() {
	cfg := internal.ParseConfig()
	logger := log.NewLogger(cfg.DebugMode)

	dbConfig := database.Config(cfg.Database)
	dbPool, err := database.Connect(&dbConfig, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}

	repos := repositories.New()
	svcs := services.New(repos)

	app := internal.NewApplication(cfg, logger, dbPool, svcs)

	// A waitgroup is established to ensure background tasks are completed before shutting down the server.
	backgroundWg := &sync.WaitGroup{}

	// Start server. The listen chan isn't used here and is buffered to 1 so the server won't be blocked.
	err = server.ServeApi(app, backgroundWg, make(chan<- int, 1))
	if err != nil {
		panic(err)
	}
}
