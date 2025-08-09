package users

import (
	"context"
	"errors"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/domain"

	"github.com/danielgtaylor/huma/v2"
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
	dbTx, commit, rollback, err := h.app.Db.BeginTx(ctx)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	createdUser, err := h.app.Services.User.Create(ctx, dbTx, domain.User{
		Email:       req.Body.Email,
		Username:    req.Body.Username,
		FirstName:   req.Body.FirstName,
		LastName:    req.Body.LastName,
		PhoneNumber: req.Body.PhoneNumber,
		Password:    req.Body.Password,
	})
	if err != nil {
		var errDuplicateData *domain.DuplicateDataError

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

	respData := responses.NewUser.FromDomainUser(createdUser)

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	return &CreateResponse{
		Body: responses.Success(respData),
	}, nil
}
