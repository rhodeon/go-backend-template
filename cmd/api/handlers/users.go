package handlers

import (
	"context"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/models"
	"time"
)

type UserRequest struct {
	Body UserRequestBody
}

type UserRequestBody struct {
	Username string `json:"username" required:"true"`
	Email    string `json:"email" required:"true" format:"email"`
}

type UserResponse struct {
	Body UserResponseBody
}

type UserResponseBody struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUser(app *internal.Application) func(context.Context, *UserRequest) (*UserResponse, error) {
	return func(ctx context.Context, req *UserRequest) (*UserResponse, error) {
		createdUser, err := app.Services.User.Create(ctx, app.DbPool, models.User{
			Username: req.Body.Username,
			Email:    req.Body.Email,
		})
		if err != nil {
			return nil, err
		}

		return &UserResponse{UserResponseBody(createdUser)}, nil
	}
}
