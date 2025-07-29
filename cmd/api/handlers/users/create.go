package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type CreateRequest struct {
	Body CreateRequestBody
}
type CreateRequestBody struct {
	Username  string `json:"username" required:"true" minLength:"1" example:"johndoe"`
	FirstName string `json:"first_name" required:"true" minLength:"1" example:"John"`
	LastName  string `json:"last_name" required:"true" minLength:"1" example:"Doe"`
	Email     string `json:"email" required:"true" format:"email" example:"johndoe@example.com"`
	Password  string `json:"password" required:"true" minLength:"1"`
	Phone     string `json:"phone" required:"false" minLength:"1"`
}

type CreateResponseData struct {
	Id        int    `json:"id" required:"true" example:"1"`
	Username  string `json:"username" required:"true" example:"johndoe"`
	FirstName string `json:"first_name" required:"true" example:"John"`
	LastName  string `json:"last_name" required:"true" example:"Doe"`
	Email     string `json:"email" required:"true" example:"johndoe@example.com"`
	Phone     string `json:"phone" required:"false"`
}

func (h *Handlers) Create(_ context.Context, _ *CreateRequest) (*responses.Envelope[CreateResponseData], error) {
	return responses.Success(CreateResponseData{}), nil
}
