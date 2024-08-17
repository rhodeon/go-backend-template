package handlers

import (
	"context"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

func Ping(ctx context.Context, input *struct{}) (*responses.PingResponse, error) {
	resp := &responses.PingResponse{}
	resp.Body.Status = "OK"
	return resp, nil
}
