package middleware

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

// Timeout establishes the duration a request can run before being terminated.
func Timeout(app *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		timeoutCtx, cancel := context.WithTimeout(ctx.Context(), app.Config.Server.RequestTimeout)
		defer cancel()

		ctx = huma.WithContext(ctx, timeoutCtx)
		next(ctx)
	}
}
