package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/go-errors/errors"
)

// ServeApi starts up a server with the app data.
// listenPort is used to optionally notify the caller that the server is available to accept connections, along with the connected port number.
// This is useful to guarantee the server is active before proceeding. e.g. for tests.
func ServeApi(ctx context.Context, app *internal.Application, backgroundWaitGroup *sync.WaitGroup, listenPort chan<- int) error {
	serverConfig := app.Config.Server
	rtr := router(app)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.HttpPort),
		IdleTimeout:  serverConfig.IdleTimeout,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		Handler:      rtr,
	}

	// Start a background goroutine to intercept and handle shutdown events.
	shutdownError := make(chan error)
	go handleShutdown(ctx, srv, shutdownError, backgroundWaitGroup, app)

	// Listen and Serve are split into 2 steps here to enable notifying the caller that the
	// server is listening for connections.
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return errors.Errorf("binding to address: %w", err)
	}

	listenPort <- listener.Addr().(*net.TCPAddr).Port //nolint: errcheck

	// Now the server can proceed to process connections.
	slog.InfoContext(
		ctx,
		"Starting server",
		slog.String(log.AttrEnv, app.Config.Environment),
		slog.String(log.AttrPort, srv.Addr),
	)

	if err := srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		// Return unsuccessful server shutdown errors.
		return errors.Errorf("shutting down server: %w", err)
	}

	// Block flow until the shutdown error channel is updated.
	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Stopped server")
	return nil
}

// handleShutdown gracefully handles interruption and termination signals,
// giving ongoing request a 20-second leeway before shutting down the server.
// It should be run as a background goroutine.
func handleShutdown(ctx context.Context, server *http.Server, shutdownErr chan error, backgroundWg *sync.WaitGroup, app *internal.Application) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait until the quit channel is updated with a signal.
	s := <-quit

	// Timeout context to delay shutdown.
	ctx, cancel := context.WithTimeout(ctx, app.Config.Server.ShutdownTimeout)
	defer cancel()

	slog.InfoContext(ctx, "Shutting down server", slog.String(log.AttrSignal, s.String()))

	// Shut down the server and update the error channel to resume execution on the main goroutine.
	err := server.Shutdown(ctx)
	if err != nil {
		shutdownErr <- err
	}

	// Wait for background tasks to complete before shutting down the application
	slog.InfoContext(ctx, "Completing background tasks...")
	backgroundWg.Wait()
	shutdownErr <- nil
}
