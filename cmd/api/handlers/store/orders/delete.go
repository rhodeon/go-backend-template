package orders

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type DeleteRequest struct {
	OrderId int `path:"order_id"`
}

type DeleteResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *Handlers) delete(_ context.Context, _ *DeleteRequest) (*DeleteResponse, error) {
	return &DeleteResponse{
		Body: responses.Success[responses.SuccessMessage]("Order deleted"),
	}, nil
}
