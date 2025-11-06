package middleware

import (
	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps handlers in an OpenTelemetry span with relevant details to trace the request.
func Tracer(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	tracer := otel.GetTracerProvider().Tracer(app.Config.Otel.ServiceName)

	return func(ctx huma.Context, next func(huma.Context)) {
		newCtx, span := tracer.Start(
			ctx.Context(),
			ctx.Operation().OperationID,
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(ctx.Method()),
				semconv.URLPath(ctx.URL().Path),
				semconv.ServerAddress(ctx.Host()),
				semconv.ServerPort(app.Config.Server.HttpPort),
				semconv.ClientAddress(ctx.RemoteAddr()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)

		defer func() {
			span.SetAttributes(semconv.HTTPResponseStatusCode(ctx.Status()))
			span.End()
		}()

		ctx = huma.WithContext(ctx, newCtx)
		next(ctx)
	}
}
