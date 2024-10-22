package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/helpers"
)

// Logger embeds a logger into the context to be accessible and modified throughout the lifetime of the request.
func Logger(app *internal.Application) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		loggerCtx := helpers.ContextWithLogger(ctx.Context(), app.Logger)
		ctx = huma.WithContext(ctx, loggerCtx)
		next(ctx)
	}
}
