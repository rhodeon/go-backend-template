package middleware

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// RequestsCounter tracks the total number of requests received.
func RequestsCounter(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	meter := otel.GetMeterProvider().Meter(app.Config.Otel.ServiceName)

	requestsCounter, err := meter.Int64Counter(
		"http.server.total_requests",
		metric.WithDescription("Total number of HTTP requests received."),
	)
	if err != nil {
		log.Fatal(context.Background(), "Failed to create requests counter", slog.Any(log.AttrError, err))
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		requestsCounter.Add(ctx.Context(), 1, metric.WithAttributes(
			semconv.HTTPRequestMethodKey.String(ctx.Method()),
			semconv.URLPath(ctx.URL().Path),
			semconv.ServerAddress(ctx.Host())))

		next(ctx)
	}
}
