package orders

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type DeleteByIdRequest struct {
	OrderId int `path:"order_id"`
}

type DeleteByIdResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *Handlers) DeleteById(_ context.Context, _ *DeleteByIdRequest) (*DeleteByIdResponse, error) {
	return &DeleteByIdResponse{
		Body: responses.Success[responses.SuccessMessage]("Order deleted"),
	}, nil
}
