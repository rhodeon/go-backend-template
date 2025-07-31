package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetByUsernameRequest struct {
	Username string `path:"username"`
}

func (h *Handlers) GetByUsername(_ context.Context, _ *struct{}) (*responses.Envelope[responses.User], error) {
	return responses.Success(responses.User{}), nil
}
