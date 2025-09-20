package otel

import (
	"context"

	"github.com/go-errors/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func SetupTracer(ctx context.Context, cfg *Config) error {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OtlpGrpcEndpoint()),
	}

	if !cfg.OtlpSecureConnection {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	otlpTraceExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
	if err != nil {
		return errors.Errorf("setting up OTLP trace exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(otlpTraceExporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.ServiceName),
			),
		),
	)

	otel.SetTracerProvider(traceProvider)
	return nil
}
