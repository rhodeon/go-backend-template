package middleware

import (
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/utils/contextutils"

	"github.com/danielgtaylor/huma/v2"
)

// SetLogger embeds a logger into the context to be accessible and modified throughout the lifetime of the request.
func SetLogger(app *internal.Application) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		loggerCtx := contextutils.WithLogger(ctx.Context(), app.Logger)
		ctx = huma.WithContext(ctx, loggerCtx)
		next(ctx)
	}
}
