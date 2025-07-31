package store

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/handlers/store/orders"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	app    *internal.Application
	Orders *orders.Handlers
}

func New(app *internal.Application, api huma.API) *Handlers {
	group := huma.NewGroup(api, "/store")

	handler := &Handlers{
		app,
		orders.New(app, group),
	}

	handler.registerRoutes(group)
	return handler
}

func (h *Handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "store-inventory",
			Method:      http.MethodGet,
			Path:        "/inventory",
			Tags:        []string{"store"},
			Description: "Returns a map of status codes to quantities.",
			Summary:     "Returns pet inventories by status",
		},
		h.Inventory,
	)
}
