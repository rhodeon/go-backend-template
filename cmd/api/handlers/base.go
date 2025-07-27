package handlers

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

// baseHandler contains the handlers for the base route `/`.
type baseHandler struct {
	app *internal.Application
}

func newBaseHandler(app *internal.Application) *baseHandler {
	return &baseHandler{app: app}
}

func (h *baseHandler) Ping(_ context.Context, _ *struct{}) (*responses.Envelope[string], error) {
	return responses.Success(responses.SuccessMessage), nil
}
