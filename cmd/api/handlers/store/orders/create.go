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

type CreateResponseData struct {
	Id       int       `json:"id" required:"true"`
	PetId    int       `json:"pet_id" required:"true"`
	Quantity int       `json:"quantity" required:"true"`
	ShipDate time.Time `json:"ship_date" required:"true"`
	Status   string    `json:"status" enum:"placed,approved,delivered" required:"true"`
	Complete bool      `json:"complete" required:"true"`
}

func (h *Handlers) Create(_ context.Context, _ *CreateRequest) (*responses.Envelope[CreateResponseData], error) {
	return responses.Success(CreateResponseData{}), nil
}
