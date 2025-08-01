package orders

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetRequest struct {
	OrderId int `path:"order_id"`
}

type GetResponse struct {
	Body responses.Envelope[responses.Order]
}

func (h *Handlers) Get(_ context.Context, _ *GetRequest) (*GetResponse, error) {
	return &GetResponse{}, nil
}
