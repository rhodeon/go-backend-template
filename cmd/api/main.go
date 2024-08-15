package main

import (
	"context"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/config"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"
)

const ConfigPrefix = "API"

func main() {
	cfg := config.Parse(ConfigPrefix)
	logger := log.NewLogger(cfg)
	app := internal.NewApplication(cfg, logger)

	dbPool, err := database.Connect(cfg, logger)
	if err != nil {
		panic(err)
	}

	res, err := dbPool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		panic(err)
	}

	app.Logger.Info(res.String())
}
