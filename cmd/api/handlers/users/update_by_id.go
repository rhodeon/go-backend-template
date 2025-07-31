package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type UpdateByIdRequest struct {
	Body UpdateByIdRequestBody
}

type UpdateByIdRequestBody struct {
	Id        int    `json:"id" required:"true" example:"1"`
	Username  string `json:"username" required:"true" example:"johndoe"`
	FirstName string `json:"first_name" required:"true" example:"John"`
	LastName  string `json:"last_name" required:"true" example:"Doe"`
	Email     string `json:"email" required:"true" example:"johndoe@example.com"`
	Phone     string `json:"phone" required:"false"`
	Password  string `json:"password" required:"true"`
}

func (h *Handlers) UpdateById(_ context.Context, _ *struct{}) (*responses.Envelope[responses.User], error) {
	return responses.Success(responses.User{}), nil
}
