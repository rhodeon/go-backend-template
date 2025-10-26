package otel

import (
	"context"

	"github.com/go-errors/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// setupMeter initialises a global meter provider to be used by both first and third-party packages.
func setupMeter(ctx context.Context, cfg *Config) error {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OtlpGrpcEndpoint()),
	}
	if !cfg.OtlpSecureConnection {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return errors.Errorf("setting up OTLP metric exporter: %w", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.ServiceName),
			),
		),
	)

	otel.SetMeterProvider(provider)
	return nil
}
