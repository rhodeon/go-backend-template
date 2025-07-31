package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type CreateRequest struct {
	Body CreateRequestBody
}

type CreateRequestBody struct {
	Username  string `json:"username" required:"true" example:"johndoe"`
	FirstName string `json:"first_name" required:"true" example:"John"`
	LastName  string `json:"last_name" required:"true" example:"Doe"`
	Email     string `json:"email" required:"true" example:"johndoe@example.com"`
	Phone     string `json:"phone" required:"false"`
	Password  string `json:"password" required:"true"`
}

func (h *Handlers) Create(_ context.Context, _ *CreateRequest) (*responses.Envelope[responses.User], error) {
	return responses.Success(responses.User{}), nil
}
