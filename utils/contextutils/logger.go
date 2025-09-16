package contextutils

import (
	"context"
	"log/slog"

	slogctx "github.com/veqryn/slog-context"
)

// withLoggerAttrs extends the given context to include the provided slog attributes.
func withLoggerAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	newCtx := ctx
	for _, attr := range attrs {
		newCtx = slogctx.Append(newCtx, attr)
	}

	return newCtx
}
