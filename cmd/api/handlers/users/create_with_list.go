package users

import (
	"context"

	handlerutils "github.com/rhodeon/go-backend-template/cmd/api/handlers/utils"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type CreateWithListRequest struct {
	Body []CreateWithListRequestBody
}
type CreateWithListRequestBody struct {
	Username  string `json:"username" required:"true" minLength:"1" example:"johndoe"`
	FirstName string `json:"first_name" required:"true" minLength:"1" example:"John"`
	LastName  string `json:"last_name" required:"true" minLength:"1" example:"Doe"`
	Email     string `json:"email" required:"true" format:"email" example:"johndoe@example.com"`
	Password  string `json:"password" required:"true" minLength:"1"`
	Phone     string `json:"phone" required:"false" minLength:"1"`
}
type CreateWithListResponseData []responses.User

func (rd CreateWithListResponseData) Name() string {
	return handlerutils.GenerateSchemaName(rd)
}

func (h *Handlers) CreateWithList(_ context.Context, _ *struct{}) (*responses.EnvelopeWithMetadata[CreateWithListResponseData, *responses.Pagination], error) {
	return responses.SuccessWithMetadata(CreateWithListResponseData{}, responses.CalculatePagination(0, 0, 0)), nil
}
