package helpers

import (
	"context"
	"log/slog"
)

type contextKey string

const loggerKey contextKey = "logger"

func ContextWithLogger(ctx context.Context, logger *slog.Logger, attrs ...slog.Attr) context.Context {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return context.WithValue(ctx, loggerKey, logger)
}

func GetContextLogger(ctx context.Context) *slog.Logger {
	if logger, exists := ctx.Value(loggerKey).(*slog.Logger); exists {
		return logger
	} else {
		return slog.Default()
	}
}
