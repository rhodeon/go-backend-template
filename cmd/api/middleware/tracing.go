package middleware

import (
	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracing wraps handlers in an OpenTelemetry span with relevant details to trace the request.
func Tracing(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		tracer := otel.GetTracerProvider().Tracer(app.Config.Otel.ServiceName)
		spanName := ctx.Operation().OperationID

		newCtx, span := tracer.Start(ctx.Context(), spanName,
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(ctx.Method()),
				semconv.URLPath(ctx.URL().Path),
				semconv.ServerAddress(ctx.Host()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		ctx = huma.WithContext(ctx, newCtx)
		next(ctx)

		span.SetAttributes(
			semconv.HTTPResponseStatusCode(ctx.Status()),
		)
	}
}
