package orders

import (
	"context"
	"time"

	handlerutils "github.com/rhodeon/go-backend-template/cmd/api/handlers/utils"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type CreateRequest struct {
	Body CreateRequestBody
}

type CreateRequestBody struct {
	PetId    int       `json:"pet_id" required:"true"`
	Quantity int       `json:"quantity" required:"true"`
	ShipDate time.Time `json:"ship_date" required:"true"`
}

func (rb CreateRequestBody) Name() string {
	return handlerutils.GenerateSchemaName(rb)
}

func (h *Handlers) Create(_ context.Context, _ *CreateRequest) (*responses.Envelope[responses.Order], error) {
	return responses.Success(responses.Order{}), nil
}
