package middleware

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rhodeon/go-backend-template/internal/helpers"
)

// Logger embeds a logger into the context to be accessible and modified throughout the lifetime of the request.
func Logger(mainContext context.Context) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		logger := helpers.ContextGetLogger(mainContext)
		loggerCtx := helpers.ContextSetLogger(ctx.Context(), logger)
		ctx = huma.WithContext(ctx, loggerCtx)
		next(ctx)
	}
}
