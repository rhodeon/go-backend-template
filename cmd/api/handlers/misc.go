package handlers

import (
	"context"
)

type PingResponse struct {
	Body struct {
		Status string `json:"status" enum:"OK" doc:"Acknowledgement status"`
	}
}

func Ping(ctx context.Context, input *struct{}) (*PingResponse, error) {
	resp := &PingResponse{}
	resp.Body.Status = "OK"
	return resp, nil
}
