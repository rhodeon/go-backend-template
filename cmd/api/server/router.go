package server

import (
	"net/http"

	"github.com/rhodeon/go-backend-template/cmd/api/handlers"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	apimiddleware "github.com/rhodeon/go-backend-template/cmd/api/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func router(app *internal.Application) http.Handler {
	mux := http.NewServeMux()

	humaConfig := newHumaConfig("API", "0.1.0")
	api := humago.New(mux, humaConfig)

	api.UseMiddleware(
		apimiddleware.RequestId(app, api),
		apimiddleware.Tracer(app, api),
		apimiddleware.Logger(app, api),
		apimiddleware.Recover(app, api),
		apimiddleware.Timeout(app, api),
	)

	api.OpenAPI().Tags = []*huma.Tag{
		{Name: "misc"},
		{Name: "user"},
		{Name: "store"},
		{Name: "pets"},
	}

	handlers.Setup(app, api)
	return mux
}
