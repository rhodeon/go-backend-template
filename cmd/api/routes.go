package main

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/rhodeon/go-backend-template/cmd/api/handlers"
	"net/http"
)

func routes() http.Handler {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("API", "0.1.0"))

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

	return router
}
