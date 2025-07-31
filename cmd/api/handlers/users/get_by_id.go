package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetByIdRequest struct {
	UserId string `path:"user_id"`
}

func (h *Handlers) GetById(_ context.Context, _ *struct{}) (*responses.Envelope[responses.User], error) {
	return responses.Success(responses.User{}), nil
}
