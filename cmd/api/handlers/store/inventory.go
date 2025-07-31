package store

import (
	"context"

	handlerutils "github.com/rhodeon/go-backend-template/cmd/api/handlers/utils"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type InventoryResponseData map[string]string

func (rd InventoryResponseData) Name() string {
	return handlerutils.GenerateSchemaName(rd)
}

func (h *Handlers) Inventory(_ context.Context, _ *struct{}) (*responses.Envelope[InventoryResponseData], error) {
	return responses.Success(InventoryResponseData{}), nil
}
