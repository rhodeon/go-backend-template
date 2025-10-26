package otel

import (
	"context"

	"github.com/go-errors/errors"
)

func SetupProviders(ctx context.Context, cfg *Config) error {
	if err := setupTracer(ctx, cfg); err != nil {
		return errors.Errorf("setting up tracer provider: %w", err)
	}

	if err := setupMeter(ctx, cfg); err != nil {
		return errors.Errorf("setting up meter provider: %w", err)
	}

	return nil
}
