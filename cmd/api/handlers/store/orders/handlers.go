package orders

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	app *internal.Application
}

func New(app *internal.Application, api huma.API) *Handlers {
	group := huma.NewGroup(api, "/orders")

	handlers := &Handlers{app: app}
	handlers.registerRoutes(group)

	return handlers
}

func (h *Handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "store-orders-create",
			Method:      http.MethodPost,
			Path:        "",
			Tags:        []string{"store"},
			Description: "Place a new order in the store.",
			Summary:     "Place an order for a pet",
		},
		h.Create,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "store-orders-get",
			Method:      http.MethodGet,
			Path:        "/{order_id}",
			Tags:        []string{"store"},
			Description: "For valid response try integer IDs with value <= 5 or > 10. Other values will generate exceptions.",
			Summary:     "Find purchase order by ID",
		},
		h.Get,
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "store-orders-delete",
			Method:      http.MethodDelete,
			Path:        "/{order_id}",
			Tags:        []string{"store"},
			Description: "For valid response try integer IDs with value < 1000. Anything above 1000 or non-integers will generate API errors.",
			Summary:     "Delete purchase order by identifier",
		},
		h.Delete,
	)
}
