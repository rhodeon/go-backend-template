package handlers

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

func (h *Handlers) Ping(_ context.Context, _ *struct{}) (*responses.Envelope[string], error) {
	return responses.Success(responses.SuccessMessage), nil
}
