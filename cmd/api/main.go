package main

import (
	"context"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
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
	app := internal.NewApplication(cfg, logger, repos, dbPool)

	res, err := app.DbPool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		panic(err)
	}
	app.Logger.Info(res.String())
}
