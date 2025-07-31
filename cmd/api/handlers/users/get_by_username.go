package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetByUsernameRequest struct {
	Username string `path:"username"`
}

type GetByUsernameResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) GetByUsername(_ context.Context, _ *struct{}) (*GetByUsernameResponse, error) {
	return &GetByUsernameResponse{}, nil
}
