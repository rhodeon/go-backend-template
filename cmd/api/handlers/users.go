package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
	api_errors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	domain_errors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"
	"time"
)

type CreateUserRequest struct {
	Body CreateUserRequestBody
}

type CreateUserRequestBody struct {
	Username string `json:"username" required:"true"`
	Email    string `json:"email" required:"true" format:"email"`
}

type GetUserRequest struct {
	Id int `json:"id" path:"id"`
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

func CreateUser(app *internal.Application) func(context.Context, *CreateUserRequest) (*UserResponse, error) {
	return func(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
		createdUser, err := app.Services.User.Create(ctx, app.DbPool, models.User{
			Username: req.Body.Username,
			Email:    req.Body.Email,
		})
		if err != nil {
			var duplicateErr *domain_errors.DuplicateDataError
			switch {
			case errors.As(err, &duplicateErr):
				return nil, huma.Error409Conflict(err.Error())

			default:
				return nil, api_errors.HandleUntypedError(ctx, err)
			}
		}

		return &UserResponse{UserResponseBody(createdUser)}, nil
	}
}

func GetUser(app *internal.Application) func(context.Context, *GetUserRequest) (*UserResponse, error) {
	return func(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
		user, err := app.Services.User.GetById(ctx, app.DbPool, int32(req.Id))
		if err != nil {
			var notFoundErr *domain_errors.RecordNotFoundError
			switch {
			case errors.As(err, &notFoundErr):
				return nil, huma.Error404NotFound(err.Error())

			default:
				return nil, api_errors.HandleUntypedError(ctx, err)
			}
		}

		return &UserResponse{UserResponseBody(user)}, nil
	}
}
