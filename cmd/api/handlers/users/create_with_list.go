package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type CreateWithListRequest struct {
	Body CreateWithListRequestBody
}

type (
	CreateWithListRequestBody []CreateWithListRequestItem
	CreateWithListRequestItem struct {
		Username  string `json:"username" required:"true" minLength:"1" example:"johndoe"`
		FirstName string `json:"first_name" required:"true" minLength:"1" example:"John"`
		LastName  string `json:"last_name" required:"true" minLength:"1" example:"Doe"`
		Email     string `json:"email" required:"true" format:"email" example:"johndoe@example.com"`
		Password  string `json:"password" required:"true" minLength:"1"`
		Phone     string `json:"phone" required:"false" minLength:"1"`
	}
)

type CreateWithListResponse struct {
	Body responses.EnvelopeWithMetadata[[]responses.User, *responses.Pagination]
}

func (h *Handlers) CreateWithList(_ context.Context, _ *CreateWithListRequest) (*CreateWithListResponse, error) {
	return &CreateWithListResponse{
		Body: responses.SuccessWithMetadata(
			[]responses.User{},
			responses.CalculatePagination(0, 0, 0),
		),
	}, nil
}
