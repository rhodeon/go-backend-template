package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type DeleteByIdRequest struct {
	UserId string `path:"user_id"`
}

func (h *Handlers) DeleteById(_ context.Context, _ *DeleteByIdRequest) (*responses.Envelope[responses.SuccessMessageResponseData], error) {
	return responses.Success(responses.SuccessMessageResponseData("User deleted")), nil
}
