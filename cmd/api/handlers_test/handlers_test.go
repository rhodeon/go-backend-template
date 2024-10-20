package handlers_test

import (
	"fmt"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/server"
	"github.com/rhodeon/go-backend-template/domain/services"
	"github.com/rhodeon/go-backend-template/repositories"
	"log/slog"
	"os"
	"sync"
	"testing"
)

var app *internal.Application
var config *internal.Config
var baseUrl string

// TestMain sets up a server and data common to all tests in the package.
// This allows the testing of routing and endpoint logic without being coupled to the web framework used.
func TestMain(m *testing.M) {
	config = &internal.Config{
		Environment: "testing",
		DebugMode:   false,
		Database:    internal.DatabaseConfig{},
		Server: internal.ServerConfig{
			HttpPort: 0, // A port of 0 makes the server connect to any available port.
		},
	}

	repos := repositories.New()
	services := services.New(repos)
	// TODO: Add database pool.
	app = internal.NewApplication(config, slog.Default(), nil, services)

	listenChan := make(chan int, 1)

	go func() {
		if err := server.ServeApi(app, &sync.WaitGroup{}, listenChan); err != nil {
			panic(err)
		}
	}()

	// Ensure the server is ready before proceeding.
	baseUrl = fmt.Sprintf("http://localhost:%d", <-listenChan)

	m.Run()
	os.Exit(0)
}
