package orders

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetByIdRequest struct {
	OrderId int `path:"order_id"`
}

type GetByIdResponse struct {
	Body responses.Envelope[responses.Order]
}

func (h *Handlers) GetById(_ context.Context, _ *GetByIdRequest) (*GetByIdResponse, error) {
	return &GetByIdResponse{}, nil
}
