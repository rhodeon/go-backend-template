package main

import (
	"context"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"
)

func main() {
	cfg := internal.ParseConfig()
	logger := log.NewLogger(cfg.DebugMode)

	dbConfig := database.Config(cfg.Database)
	dbPool, err := database.Connect(&dbConfig, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}

	res, err := dbPool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		panic(err)
	}

	app := internal.NewApplication(cfg, logger)
	app.Logger.Info(res.String())
}
