package helpers

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	contextKeylogger    contextKey = "logger"
	ContextKeyRequestId contextKey = "request_id"
)

func ContextWithLogger(ctx context.Context, logger *slog.Logger, attrs ...slog.Attr) context.Context {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return context.WithValue(ctx, contextKeylogger, logger)
}

func GetContextLogger(ctx context.Context) *slog.Logger {
	if logger, exists := ctx.Value(contextKeylogger).(*slog.Logger); exists {
		return logger
	} else {
		return slog.Default()
	}
}

func UpdateContextLogger(ctx context.Context, attrs ...slog.Attr) context.Context {
	logger := GetContextLogger(ctx)
	return ContextWithLogger(ctx, logger, attrs...)
}
