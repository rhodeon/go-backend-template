package auth

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	app *internal.Application
}

func New(app *internal.Application, api huma.API) *Handlers {
	group := huma.NewGroup(api, "/auth")

	handlers := &Handlers{app: app}
	handlers.registerRoutes(group)
	return handlers
}

func (h *Handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "auth-register",
			Method:      http.MethodPost,
			Path:        "/register",
			Tags:        []string{"auth"},
			Summary:     "Register user",
			Description: "Registers a new user. Email verification is required before further actions can be performed.",
		},
		h.register,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "auth-verify-account",
			Method:      http.MethodPost,
			Path:        "/verify-account",
			Tags:        []string{"auth"},
			Summary:     "Verify account",
			Description: "Activates the user as a verified account.",
		},
		h.verifyAccount,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "auth-send-verification-email",
			Method:      http.MethodPost,
			Path:        "/verification-email",
			Tags:        []string{"auth"},
			Summary:     "Send verification email",
			Description: "Sends verification email to unverified users.",
		},
		h.sendVerificationEmail,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "auth-login",
			Method:      http.MethodPost,
			Path:        "/login",
			Tags:        []string{"auth"},
			Summary:     "Login user",
			Description: "Log into the system.",
		},
		h.login,
	)
}
