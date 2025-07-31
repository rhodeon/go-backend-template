package users

import (
	"context"

	handlerutils "github.com/rhodeon/go-backend-template/cmd/api/handlers/utils"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type LoginRequest struct {
	Body LoginRequestBody
}

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (rb LoginRequestBody) Name() string {
	return handlerutils.GenerateSchemaName(rb)
}

// TODO: Add response headers.
func (h *Handlers) Login(_ context.Context, _ *struct{}) (*responses.Envelope[responses.SuccessMessageResponseData], error) {
	return responses.Success(responses.SuccessMessageResponseData("Success")), nil
}
