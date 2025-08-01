package users

import (
	"context"
	"time"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type LoginRequest struct {
	Body LoginRequestBody
}

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Body          responses.Envelope[responses.SuccessMessage]
	XRateLimit    int    `header:"X-Rate-Limit" doc:"Calls per hour allowed by the user."`
	XExpiresAfter string `header:"X-Expires-After" doc:"Date in UTC when token expires."`
}

func (h *Handlers) login(_ context.Context, _ *struct{}) (*LoginResponse, error) {
	return &LoginResponse{
		Body:          responses.Success(responses.SuccessMessage("Success")),
		XRateLimit:    10,
		XExpiresAfter: time.Now().UTC().String(),
	}, nil
}
