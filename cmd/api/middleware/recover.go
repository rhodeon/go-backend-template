package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
	api_errors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"net/http"
)

// Recover writes an internal server error to the response.
func Recover(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			if err := recover(); err != nil {
				_ = huma.WriteErr(
					api,
					ctx,
					http.StatusInternalServerError,
					"",

					// Two errors are written into the response. The first contains the actual error,
					// and the second is to signify to api_errors.NewApiError() that the source of this is from a panic during logging.
					errors.Errorf("%v", err),
					api_errors.ErrPanic,
				)
			}
		}()

		next(ctx)
	}
}
