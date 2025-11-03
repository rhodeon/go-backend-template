package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// ErrorCounter tracks the total number of client and server errors.
func ErrorCounter(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	meter := otel.GetMeterProvider().Meter(app.Config.Otel.ServiceName)

	clientErrorsCounter, err := meter.Int64Counter(
		"http.server.client_error_count",
		metric.WithDescription("Total number of client errors."),
	)
	if err != nil {
		log.Fatal(context.Background(), "Failed to create client errors counter", slog.Any(log.AttrError, err))
	}

	serverErrorsCounter, err := meter.Int64Counter(
		"http.server.server_error_count",
		metric.WithDescription("Total number of server errors."))
	if err != nil {
		log.Fatal(context.Background(), "Failed to create server errors counter", slog.Any(log.AttrError, err))
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			switch {
			case isClientError(ctx.Status()):
				clientErrorsCounter.Add(ctx.Context(), 1)

			case isServerError(ctx.Status()):
				serverErrorsCounter.Add(ctx.Context(), 1)
			}
		}()

		next(ctx)
	}
}

// isClientError returns true if the status is in the 4XX range.
func isClientError(status int) bool {
	return http.StatusBadRequest <= status && status <= http.StatusUnavailableForLegalReasons
}

// isClientError returns true if the status is in the 5XX range.
func isServerError(status int) bool {
	return http.StatusInternalServerError <= status && status <= http.StatusNetworkAuthenticationRequired
}
