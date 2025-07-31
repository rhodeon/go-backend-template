package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetByIdRequest struct {
	UserId int `path:"user_id"`
}

type GetByIdResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) GetById(_ context.Context, _ *struct{}) (*GetByIdResponse, error) {
	return &GetByIdResponse{
		Body: responses.Success(responses.User{}),
	}, nil
}
