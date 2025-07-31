package users

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	app *internal.Application
}

func New(app *internal.Application, api huma.API) *Handlers {
	group := huma.NewGroup(api, "/users")

	handlers := &Handlers{app: app}
	handlers.registerRoutes(group)
	return handlers
}

func (h *Handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-create",
			Method:      http.MethodPost,
			Path:        "",
			Tags:        []string{"users"},
			Description: "This can only be done by the logged in user.",
			Summary:     "Create user",
		},
		h.Create,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-create-with-list",
			Method:      http.MethodPost,
			Path:        "/create-with-list",
			Tags:        []string{"users"},
			Summary:     "Creates list of users with given input array",
			Description: "Creates list of users with given input array.",
		},
		h.CreateWithList,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-login",
			Method:      http.MethodPost,
			Path:        "/login",
			Tags:        []string{"users"},
			Summary:     "Logs user into the system",
			Description: "Log into the system.",
		},
		h.Login,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-logout",
			Method:      http.MethodPost,
			Path:        "/logout",
			Tags:        []string{"users"},
			Summary:     "Logs out current logged in user session",
			Description: "Log user out of the system.",
		},
		h.Logout,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-get-by-id",
			Method:      http.MethodGet,
			Path:        "/{user_id}",
			Tags:        []string{"users"},
			Summary:     "Get user by id",
			Description: "Get user detail based on id.",
		},
		h.GetById,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-get-by-username",
			Method:      http.MethodGet,
			Path:        "/{username}",
			Tags:        []string{"users"},
			Summary:     "Get user by username",
			Description: "Get user detail based on username.",
		},
		h.GetByUsername,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-update-by-id",
			Method:      http.MethodPut,
			Path:        "/{user_id}",
			Tags:        []string{"users"},
			Summary:     "Update user resource",
			Description: "This can only be done by the logged in user.",
		},
		h.UpdateById,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-delete-by-id",
			Method:      http.MethodPost,
			Path:        "/{id}",
			Tags:        []string{"users"},
			Summary:     "Delete user resource",
			Description: "This can only be done by the logged in user.",
		},
		h.DeleteById,
	)
}
