package users

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	domainerrors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"
	"github.com/rhodeon/go-backend-template/internal/database"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
)

type CreateRequest struct {
	Body CreateRequestBody
}

type CreateRequestBody struct {
	Email       string `json:"email" required:"true" example:"johndoe@example.com" minLength:"3"`
	Username    string `json:"username" required:"true" example:"johndoe" minLength:"3"`
	FirstName   string `json:"first_name" required:"true" example:"John" minLength:"1"`
	LastName    string `json:"last_name" required:"true" example:"Doe" minLength:"1"`
	PhoneNumber string `json:"phone_number" required:"false" minLength:"1"`
	Password    string `json:"password" required:"true" example:"password" minLength:"6"`
}

type CreateResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	dbTx, commit, rollback, err := database.BeginTransaction(ctx, h.app.DbPool)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	createdUser, err := h.app.Services.User.Create(ctx, dbTx, models.User{
		Email:       req.Body.Email,
		Username:    req.Body.Username,
		FirstName:   req.Body.FirstName,
		LastName:    req.Body.LastName,
		PhoneNumber: req.Body.PhoneNumber,
		Password:    req.Body.Password,
	})
	if err != nil {
		var errDuplicateData *domainerrors.DuplicateDataError

		switch {
		case errors.As(err, &errDuplicateData):
			switch errDuplicateData.Field {
			case "email":
				return nil, huma.Error409Conflict("email already taken")

			case "username":
				return nil, huma.Error409Conflict("username already taken")
			}

			fallthrough

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	respData := responses.User{
		Id:        createdUser.Id,
		Username:  createdUser.Username,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
		Phone:     createdUser.PhoneNumber,
		CreatedAt: createdUser.CreatedAt,
	}

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	return &CreateResponse{
		Body: responses.Success(respData),
	}, nil
}
