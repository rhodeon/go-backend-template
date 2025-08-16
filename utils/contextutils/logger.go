package contextutils

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/log"
)

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// GetLogger retrieves the logger saved in the given context. If none is found, it falls back to the default logger.
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, exists := ctx.Value(contextKeyLogger).(*slog.Logger); exists {
		return logger
	} else {
		return log.NewLogger(false)
	}
}

func WithLoggerAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	logger := GetLogger(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return WithLogger(ctx, logger)
}
