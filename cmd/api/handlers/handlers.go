package handlers

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/handlers/pets"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers/store"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers/users"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	app   *internal.Application
	Store *store.Handlers
	Users *users.Handlers
	Pets  *pets.Handlers
}

func Setup(app *internal.Application, api huma.API) {
	handlers := &Handlers{
		app,
		store.New(app, api),
		users.New(app, api),
		pets.New(app, api),
	}

	handlers.registerRoutes(api)
}

func (h *Handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "ping",
			Method:      http.MethodGet,
			Path:        "/ping",
			Tags:        []string{"misc"},
			Description: "Acknowledges that the server is reachable.",
		},
		h.Ping,
	)
}
