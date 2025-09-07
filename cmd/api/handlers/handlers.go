package handlers

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/handlers/auth"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers/pets"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers/store"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers/users"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"

	"github.com/danielgtaylor/huma/v2"
)

type handlers struct {
	app   *internal.Application
	auth  *auth.Handlers
	store *store.Handlers
	users *users.Handlers
	pets  *pets.Handlers
}

func Setup(app *internal.Application, api huma.API) {
	h := &handlers{
		app,
		auth.New(app, api),
		store.New(app, api),
		users.New(app, api),
		pets.New(app, api),
	}

	h.registerRoutes(api)
}

func (h *handlers) registerRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "ping",
			Method:      http.MethodGet,
			Path:        "/ping",
			Tags:        []string{"misc"},
			Description: "Acknowledges that the server is reachable.",
			Security:    []map[string][]string{},
		},
		h.ping,
	)
}
