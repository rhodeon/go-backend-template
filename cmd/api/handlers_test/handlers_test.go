package handlers_test

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/server"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/cache/redis"
	mockemail "github.com/rhodeon/go-backend-template/repositories/email/mock"
	"github.com/rhodeon/go-backend-template/services"
	"github.com/rhodeon/go-backend-template/utils/testutils"

	"github.com/go-errors/errors"
)

// spawnServer sets up a server and data common to all tests in the package.
// This allows the testing of routing and endpoint logic without being coupled to the web framework used.
func spawnServer() (*internal.Application, error) {
	config := &internal.Config{
		Environment: "testing",
		DebugMode:   false,
		Database:    internal.DatabaseConfig{},
		Server: internal.ServerConfig{
			HttpPort:        0, // A port of 0 makes the server connect to any available port.
			IdleTimeout:     1 * time.Minute,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			RequestTimeout:  10 * time.Second,
		},
	}

	dbPool, err := testutils.ConnectDb(context.Background())
	if err != nil {
		return nil, errors.Errorf("connecting database: %w", err)
	}

	redisCache, err := redis.NewTestCache(context.Background())
	if err != nil {
		return nil, errors.Errorf("connecting redis: %w", err)
	}

	repos := repositories.New(redisCache, mockemail.New())
	svcs := services.New(repos, &services.Config{})
	app := internal.NewApplication(config, slog.Default(), dbPool, svcs)
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
