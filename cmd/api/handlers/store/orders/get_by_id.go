package orders

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetByIdRequest struct {
	OrderId string `path:"order_id"`
}

func (h *Handlers) GetById(_ context.Context, _ *GetByIdRequest) (*responses.Envelope[responses.Order], error) {
	return responses.Success(responses.Order{}), nil
}
