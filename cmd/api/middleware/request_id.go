package middleware

import (
	"log/slog"
	"strings"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/utils/contextutils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// SetRequestId adds a unique ID to identify and filter the request logs and metrics.
func SetRequestId(_ *internal.Application) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		requestId := strings.ReplaceAll(uuid.New().String(), "-", "")
		newCtx := contextutils.WithLoggerAttrs(
			contextutils.WithRequestId(ctx.Context(), requestId),
			slog.String(log.AttrRequestId, requestId),
		)

		ctx = huma.WithContext(ctx, newCtx)
		next(ctx)
	}
}
