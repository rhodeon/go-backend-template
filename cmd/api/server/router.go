package server

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/handlers"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	apimiddleware "github.com/rhodeon/go-backend-template/cmd/api/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func router(app *internal.Application) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	humaConfig := newHumaConfig("API", "0.1.0")
	api := humachi.New(mux, humaConfig)

	api.UseMiddleware(
		apimiddleware.SetLogger(app),
		apimiddleware.SetRequestId(app),
		apimiddleware.Timeout(app),
		apimiddleware.Recover(api),
	)

	api.OpenAPI().Tags = []*huma.Tag{
		{Name: "misc"},
		{Name: "user"},
		{Name: "store"},
	}

	handlers.Setup(app, api)
	return mux
}
