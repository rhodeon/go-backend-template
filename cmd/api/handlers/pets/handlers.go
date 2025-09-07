package pets

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	apimiddleware "github.com/rhodeon/go-backend-template/cmd/api/middleware"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	app *internal.Application
}

func New(app *internal.Application, api huma.API) *Handlers {
	group := huma.NewGroup(api, "/pets")
	group.UseMiddleware(apimiddleware.Authentication(app, api))

	handler := &Handlers{
		app,
	}

	handler.registerRoutes(group)
	return handler
}

func (h *Handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "pets-create",
			Method:      http.MethodPost,
			Path:        "",
			Tags:        []string{"pets"},
			Summary:     "Create pet",
			Description: "This can only be done by the logged in user.",
		},
		h.create,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "pets-get",
			Method:      http.MethodGet,
			Path:        "/{pet_id}",
			Tags:        []string{"pets"},
			Summary:     "Get pet by id",
			Description: "Get pet detail based on id.",
		},
		h.get,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "pets-update",
			Method:      http.MethodPatch,
			Path:        "/{pet_id}",
			Tags:        []string{"pets"},
			Summary:     "Update pet resource",
			Description: "This can only be done by the logged in user.",
		},
		h.update,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "users-delete",
			Method:      http.MethodPost,
			Path:        "/{pet_id}",
			Tags:        []string{"pets"},
			Summary:     "Delete pet resource",
		},
		h.delete,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "pets-list",
			Method:      http.MethodGet,
			Path:        "",
			Tags:        []string{"pets"},
			Summary:     "List pets",
		},
		h.list,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "pets-upload-image",
			Method:      http.MethodPost,
			Path:        "/upload-image",
			Tags:        []string{"pets"},
			Summary:     "Upload pet image",
		},
		h.uploadImage,
	)
}
