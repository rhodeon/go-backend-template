package middleware

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
)

// Timeout establishes the duration a request can run before being terminated.
func Timeout(app *internal.Application) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		timeoutCtx, cancel := context.WithTimeout(ctx.Context(), app.Config.Server.RequestTimeout)
		defer cancel()

		ctx = huma.WithContext(ctx, timeoutCtx)
		next(ctx)
	}
}
