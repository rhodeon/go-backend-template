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
				// sloglint is disabled here because of the `static-msg` rule.
				// The stacktrace is dynamic and isn't set as a log attribute as newlines are not rendered.
				helpers.ContextGetLogger(ctx.Context()).Error(fmt.Sprintf( //nolint: sloglint
					"Internal server error caused by panic: %v\n%v\n",
					r, string(debug.Stack()),
				))

				_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "")
			}
		}()

		next(ctx)
	}
}
