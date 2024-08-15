package main

import (
	"context"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
)

func main() {
	cfg := internal.ParseConfig()
	logger := internal.SetupLogger(cfg)
	app := internal.NewApplication(cfg, logger)

	dbPool, err := internal.ConnectToDb(cfg, logger)
	if err != nil {
		panic(err)
	}

	res, err := dbPool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		panic(err)
	}

	app.Logger.Info(res.String())
}
