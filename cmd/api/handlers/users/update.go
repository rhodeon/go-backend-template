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
	Username  string `json:"username" required:"true" example:"johndoe"`
	FirstName string `json:"first_name" required:"true" example:"John"`
	LastName  string `json:"last_name" required:"true" example:"Doe"`
	Email     string `json:"email" required:"true" example:"johndoe@example.com"`
	Phone     string `json:"phone" required:"false"`
	Password  string `json:"password" required:"true"`
}
type UpdateResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) Update(_ context.Context, _ *UpdateRequest) (*UpdateResponse, error) {
	return &UpdateResponse{}, nil
}
