package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

func (h *Handlers) Logout(_ context.Context, _ *struct{}) (*responses.Envelope[responses.SuccessMessageResponseData], error) {
	return responses.Success(responses.SuccessMessageResponseData("Success")), nil
}
