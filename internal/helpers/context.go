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

func ContextSetLogger(ctx context.Context, logger *slog.Logger, attrs ...slog.Attr) context.Context {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return context.WithValue(ctx, contextKeylogger, logger)
}

// ContextGetLogger returns any embedded logger found in the given context, otherwise it falls back to the default.
func ContextGetLogger(ctx context.Context) *slog.Logger {
	if logger, exists := ctx.Value(contextKeylogger).(*slog.Logger); exists {
		return logger
	} else {
		return slog.Default()
	}
}

func ContextUpdateLogger(ctx context.Context, attrs ...slog.Attr) context.Context {
	logger := ContextGetLogger(ctx)
	return ContextSetLogger(ctx, logger, attrs...)
}
