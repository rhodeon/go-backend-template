package store

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type InventoryResponse struct {
	Body responses.Envelope[map[string]string]
}

func (h *Handlers) inventory(_ context.Context, _ *struct{}) (*InventoryResponse, error) {
	return &InventoryResponse{
		Body: responses.Success(map[string]string{}),
	}, nil
}
