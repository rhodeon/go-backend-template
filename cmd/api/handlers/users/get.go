package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetRequest struct {
	UserId int `path:"user_id"`
}

type GetResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) Get(_ context.Context, _ *GetRequest) (*GetResponse, error) {
	return &GetResponse{
		Body: responses.Success(responses.User{}),
	}, nil
}
