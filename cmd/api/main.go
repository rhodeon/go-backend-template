package main

import (
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
)

func main() {
	cfg := internal.ParseConfig()
	logger := internal.SetupLogger(cfg)
	app := internal.NewApplication(cfg, logger)

	app.Logger.Info("Testing 123...")
}
