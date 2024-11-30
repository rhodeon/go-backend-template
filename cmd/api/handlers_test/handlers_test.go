package handlers_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/server"
	"github.com/rhodeon/go-backend-template/domain/services"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/test_utils"
	"log/slog"
	"sync"
)

// spawnServer sets up a server and data common to all tests in the package.
// This allows the testing of routing and endpoint logic without being coupled to the web framework used.
func spawnServer() (*internal.Application, error) {
	config := &internal.Config{
		Environment: "testing",
		DebugMode:   false,
		Database:    internal.DatabaseConfig{},
		Server: internal.ServerConfig{
			HttpPort: 0, // A port of 0 makes the server connect to any available port.
		},
	}

	repos := repositories.New()
	services := services.New(repos)
	dbPool, err := test_utils.ConnectDb(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	app := internal.NewApplication(config, slog.Default(), dbPool, services)
	listenChan := make(chan int, 1)

	go func() {
		if err := server.ServeApi(app, &sync.WaitGroup{}, listenChan); err != nil {
			panic(err)
		}
	}()

	// Ensure the server is ready before proceeding.
	app.Config.Server.BaseUrl = fmt.Sprintf("http://localhost:%d", <-listenChan)

	return app, nil
}
