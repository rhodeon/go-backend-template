package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type UpdateRequest struct {
	Body   UpdateRequestBody
	UserId int `path:"user_id"`
}

type UpdateRequestBody struct {
	Username  string `json:"username" required:"false"`
	FirstName string `json:"first_name" required:"false"`
	LastName  string `json:"last_name" required:"false"`
	Email     string `json:"email" required:"false"`
	Phone     string `json:"phone" required:"false"`
	Password  string `json:"password" required:"false"`
}
type UpdateResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) update(_ context.Context, _ *UpdateRequest) (*UpdateResponse, error) {
	return &UpdateResponse{}, nil
}
