package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
	api_errors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/models/requests"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	domain_errors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"
)

func CreateUser(app *internal.Application) func(context.Context, *requests.CreateUserRequest) (*responses.UserResponse, error) {
	return func(ctx context.Context, req *requests.CreateUserRequest) (*responses.UserResponse, error) {
		createdUser, err := app.Services.User.Create(ctx, app.DbPool, models.User{
			Username: req.Body.Username,
			Email:    req.Body.Email,
		})
		if err != nil {
			//helpers.GetContextLogger(ctx).ErrorContext(ctx, "blah", slog.Any(log.AttrError, err))
			var duplicateErr *domain_errors.DuplicateDataError
			switch {
			case errors.As(err, &duplicateErr):
				return nil, huma.Error409Conflict(err.Error())

			default:
				return nil, api_errors.HandleUntypedError(ctx, err)
			}
		}

		return &responses.UserResponse{Body: responses.UserResponseBody(createdUser)}, nil
	}
}

func GetUser(app *internal.Application) func(context.Context, *requests.GetUserRequest) (*responses.UserResponse, error) {
	return func(ctx context.Context, req *requests.GetUserRequest) (*responses.UserResponse, error) {
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

		return &responses.UserResponse{Body: responses.UserResponseBody(user)}, nil
	}
}
