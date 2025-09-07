package middleware

import (
	"net/http"
	"strings"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/domain"
	"github.com/rhodeon/go-backend-template/services"
	"github.com/rhodeon/go-backend-template/utils/contextutils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-errors/errors"
)

// Authentication ensures the received auth token is valid, and stores the user ID for further use in the session.
func Authentication(app *internal.Application, api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		dbTx, commit, rollback, err := app.Db.BeginTx(ctx.Context())
		if err != nil {
			// Panicking here triggers the recover middleware which results in a 500 with the error being logged.
			panic(err)
		}
		defer rollback(ctx.Context())

		token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
		if len(token) == 0 {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthenticated")
			return
		}

		tokenClaims, err := app.Services.Auth.ParseToken(token, services.AuthTokenTypeAccess)
		if err != nil || tokenClaims.Type != string(services.AuthTokenTypeAccess) {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthenticated")
			return
		}

		// Only verified users can proceed.
		userId := tokenClaims.UserId
		user, err := app.Services.User.GetById(ctx.Context(), dbTx, userId)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrUserNotFound):
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthenticated")
				return

			default:
				panic(err)
			}
		}

		if !user.IsVerified {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "account is unverified")
			return
		}

		if err := commit(ctx.Context()); err != nil {
			panic(err)
		}

		authCtx := contextutils.WithUserId(ctx.Context(), tokenClaims.UserId)
		ctx = huma.WithContext(ctx, authCtx)

		next(ctx)
	}
}
