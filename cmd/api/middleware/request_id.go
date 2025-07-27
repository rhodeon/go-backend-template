package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/helpers"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// SetRequestId adds a unique ID to identify and filter the request logs and metrics.
func SetRequestId(_ *internal.Application) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		requestId := strings.ReplaceAll(uuid.New().String(), "-", "")
		newCtx := context.WithValue(ctx.Context(), helpers.ContextKeyRequestId, requestId)
		newCtx = helpers.UpdateContextLogger(newCtx, slog.String(log.AttrRequestId, requestId))

		ctx = huma.WithContext(ctx, newCtx)
		next(ctx)
	}
}
