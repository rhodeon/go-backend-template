package contextutils

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/log"
)

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return withLoggerAttrs(
		context.WithValue(ctx, contextKeyRequestId, requestId),
		slog.String(log.AttrRequestId, requestId),
	)
}
