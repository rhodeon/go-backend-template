package handlers

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

func Ping(_ context.Context, _ *struct{}) (*responses.PingResponse, error) {
	resp := &responses.PingResponse{
		Body: responses.PingResponseBody{
			Status: "OK",
		},
	}

	return resp, nil
}
