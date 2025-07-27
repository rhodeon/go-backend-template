package handlers

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/models/requests"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	domainerrors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
)

type usersHandler struct {
	app *internal.Application
}

func newUsersHandler(app *internal.Application) *usersHandler {
	return &usersHandler{app: app}
}

func (h *usersHandler) Create(ctx context.Context, req *requests.UsersCreateRequest) (*responses.Envelope[responses.User], error) {
	createdUser, err := h.app.Services.User.Create(ctx, h.app.DbPool, models.User{
		Username: req.Body.Username,
		Email:    req.Body.Email,
	})
	if err != nil {
		var duplicateErr *domainerrors.DuplicateDataError

		switch {
		case errors.As(err, &duplicateErr):
			return nil, huma.Error409Conflict(err.Error())

		default:
			return nil, apierrors.HandleUntypedError(ctx, err)
		}
	}

	responseData := responses.User(createdUser)
	return responses.Success(responseData), nil
}

func (h *usersHandler) GetById(ctx context.Context, req *requests.UsersGetByIdRequest) (*responses.Envelope[responses.User], error) {
	user, err := h.app.Services.User.GetById(ctx, h.app.DbPool, int32(req.Id))
	if err != nil {
		var notFoundErr *domainerrors.RecordNotFoundError

		switch {
		case errors.As(err, &notFoundErr):
			return nil, huma.Error404NotFound(err.Error())

		default:
			return nil, apierrors.HandleUntypedError(ctx, err)
		}
	}

	responseData := responses.User(user)
	return responses.Success(responseData), nil
}
