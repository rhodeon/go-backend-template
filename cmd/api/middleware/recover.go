package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rhodeon/go-backend-template/internal/helpers"

	"github.com/danielgtaylor/huma/v2"
)

// Recover logs any existing panics then writes an internal server error to the response.
func Recover(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			if r := recover(); r != nil {
				logger := helpers.GetContextLogger(ctx.Context())
				logger.ErrorContext(ctx.Context(), fmt.Sprintf("Internal server error caused by panic: %v\n%v\n", r, string(debug.Stack())))

				_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "")
			}
		}()

		next(ctx)
	}
}
