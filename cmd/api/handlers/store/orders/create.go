package orders

import (
	"context"
	"time"

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

type CreateResponse struct {
	Body responses.Envelope[responses.Order]
}

func (h *Handlers) Create(_ context.Context, _ *CreateRequest) (*CreateResponse, error) {
	return &CreateResponse{}, nil
}
