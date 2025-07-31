package orders

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type DeleteByIdRequest struct {
	OrderId string `path:"order_id"`
}

func (h *Handlers) DeleteById(_ context.Context, _ *DeleteByIdRequest) (*responses.Envelope[responses.SuccessMessageResponseData], error) {
	return responses.Success(responses.SuccessMessageResponseData("Order deleted")), nil
}
