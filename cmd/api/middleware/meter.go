package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Meter combines all generic request metrics in a single middleware for simplicity.
func Meter(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	ctx := context.Background()
	meter := otel.GetMeterProvider().Meter(app.Config.Otel.ServiceName)

	requestsDurationHistogram, err := meter.Int64Histogram(
		"http.server.requests_duration",
		metric.WithDescription("Duration of HTTP requests received."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatal(ctx, "Failed to create requests duration histogram", slog.Any(log.AttrError, err))
	}

	requestsCounter, err := meter.Int64Counter(
		"http.server.total_requests",
		metric.WithDescription("Total number of HTTP requests received."),
	)
	if err != nil {
		log.Fatal(ctx, "Failed to create total requests counter", slog.Any(log.AttrError, err))
	}

	var activeRequests atomic.Int64
	_, err = meter.Int64ObservableGauge(
		"http.server.active_requests",
		metric.WithDescription("Total number of HTTP requests currently active."),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(activeRequests.Load())
			return nil
		}),
	)
	if err != nil {
		log.Fatal(ctx, "Failed to create active requests gauge", slog.Any(log.AttrError, err))
	}

	clientErrorsCounter, err := meter.Int64Counter(
		"http.server.client_error_count",
		metric.WithDescription("Total number of returned client errors."),
	)
	if err != nil {
		log.Fatal(ctx, "Failed to create client errors counter", slog.Any(log.AttrError, err))
	}

	serverErrorsCounter, err := meter.Int64Counter(
		"http.server.server_error_count",
		metric.WithDescription("Total number of returned server errors."))
	if err != nil {
		log.Fatal(ctx, "Failed to create server errors counter", slog.Any(log.AttrError, err))
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		startTime := time.Now()

		attributeSet := attribute.NewSet(
			semconv.HTTPRequestMethodKey.String(ctx.Method()),
			semconv.URLPath(ctx.URL().Path),
			semconv.ServerAddress(ctx.Host()),
		)

		activeRequests.Add(1)
		requestsCounter.Add(ctx.Context(), 1, metric.WithAttributeSet(attributeSet))

		defer func() {
			switch {
			case isClientError(ctx.Status()):
				clientErrorsCounter.Add(ctx.Context(), 1)

			case isServerError(ctx.Status()):
				serverErrorsCounter.Add(ctx.Context(), 1)
			}

			activeRequests.Add(-1)

			requestsDurationHistogram.Record(
				ctx.Context(),
				time.Since(startTime).Milliseconds(),
				metric.WithAttributeSet(attributeSet),
			)
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
