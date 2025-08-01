package users

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type LogoutResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *Handlers) logout(_ context.Context, _ *struct{}) (*LogoutResponse, error) {
	return &LogoutResponse{
		Body: responses.Success(responses.SuccessMessage("Success")),
	}, nil
}
