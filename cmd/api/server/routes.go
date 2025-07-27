package server

import (
	"net/http"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	apimiddleware "github.com/rhodeon/go-backend-template/cmd/api/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *internal.Application) http.Handler {
	router := chi.NewMux()
	router.Use(middleware.Logger)

	// The default error structure of huma is overwritten by a custom ApiError.
	huma.NewError = apierrors.NewApiError()

	api := humachi.New(router, huma.DefaultConfig("API", "0.1.0"))
	api.UseMiddleware(
		apimiddleware.Logger(app),
		apimiddleware.SetRequestId(app),
		apimiddleware.Timeout(app),
		apimiddleware.Recover(api),
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "ping",
			Method:      http.MethodGet,
			Path:        "/ping",
			Tags:        []string{"misc"},
			Description: "Acknowledges that the server is reachable.",
		},
		handlers.Ping,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "create-user",
			Method:      http.MethodPost,
			Path:        "/users",
			Tags:        []string{"users"},
			Description: "Creates a new user.",
		},
		handlers.CreateUser(app),
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "get-user-by-id",
			Method:      http.MethodGet,
			Path:        "/users/{id}",
			Tags:        []string{"users"},
			Description: "Returns a user with the given id.",
		},
		handlers.GetUser(app),
	)

	return router
}
