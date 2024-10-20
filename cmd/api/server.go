package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// serveApi starts up a server with the app data.
func serveApi(app *internal.Application, backgroundWaitGroup *sync.WaitGroup) error {
	serverConfig := app.Config.Server
	router := routes(app)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.HttpPort),
		IdleTimeout:  serverConfig.IdleTimeout,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		Handler:      router,
	}

	// Start a background goroutine to intercept and handle shutdown events.
	shutdownError := make(chan error)
	go handleShutdown(srv, shutdownError, backgroundWaitGroup, app)

	// Start and listen on server until an error occurs.
	app.Logger.Info(
		"starting server",
		slog.String("env", app.Config.Environment),
		slog.String("port", srv.Addr),
	)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Return unsuccessful server shutdown errors.
		return errors.Wrap(err, "server shutdown")
	}

	// Block flow until the shutdown error channel is updated.
	err := <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.Info("Stopped server")
	return nil
}

// handleShutdown gracefully handles interruption and termination signals,
// giving ongoing request a 20-second leeway before shutting down the server.
// It should be run as a background goroutine.
func handleShutdown(server *http.Server, shutdownErr chan error, backgroundWg *sync.WaitGroup, app *internal.Application) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait until the quit channel is updated with a signal.
	s := <-quit

	// 20-second timeout context to delay shutdown
	ctx, cancel := context.WithTimeout(context.Background(), app.Config.Server.ShutdownTimeout)
	defer cancel()

	app.Logger.Info("Shutting down server", slog.String(log.AttrSignal, s.String()))

	// Shut down the server and update the error channel  to resume execution on the main goroutine.
	err := server.Shutdown(ctx)
	if err != nil {
		shutdownErr <- err
	}

	// Wait for background tasks to complete before shutting down the application
	app.Logger.Info("Completing background tasks...")
	backgroundWg.Wait()
	shutdownErr <- nil
}
