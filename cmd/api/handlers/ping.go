package handlers

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type PingResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *handlers) ping(_ context.Context, _ *struct{}) (*PingResponse, error) {
	return &PingResponse{
		Body: responses.Success[responses.SuccessMessage]("Success"),
	}, nil
}
