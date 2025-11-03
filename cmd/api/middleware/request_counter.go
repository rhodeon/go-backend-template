package middleware

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// RequestCounter tracks the active and total number of requests received.
func RequestCounter(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	meter := otel.GetMeterProvider().Meter(app.Config.Otel.ServiceName)

	requestsDurationHistogram, err := meter.Int64Histogram(
		"http.server.requests_duration",
		metric.WithDescription("Duration of HTTP requests received"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatal(context.Background(), "Failed to create requests duration histogram", slog.Any(log.AttrError, err))
	}

	requestsCounter, err := meter.Int64Counter(
		"http.server.total_requests",
		metric.WithDescription("Total number of HTTP requests received."),
	)
	if err != nil {
		log.Fatal(context.Background(), "Failed to create total requests counter", slog.Any(log.AttrError, err))
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
		log.Fatal(context.Background(), "Failed to create active requests gauge", slog.Any(log.AttrError, err))
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		startTime := time.Now()
		activeRequests.Add(1)

		defer func() {
			activeRequests.Add(-1)

			requestsDurationHistogram.Record(
				ctx.Context(),
				time.Since(startTime).Milliseconds(),
				metric.WithAttributes(
					semconv.HTTPRequestMethodKey.String(ctx.Method()),
					semconv.URLPath(ctx.URL().Path),
					semconv.ServerAddress(ctx.Host()),
				))
		}()

		requestsCounter.Add(ctx.Context(), 1, metric.WithAttributes(
			semconv.HTTPRequestMethodKey.String(ctx.Method()),
			semconv.URLPath(ctx.URL().Path),
			semconv.ServerAddress(ctx.Host()),
		))

		next(ctx)
	}
}
