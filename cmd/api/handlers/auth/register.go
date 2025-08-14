package auth

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-errors/errors"
)

type RegisterRequest struct {
	Body RegisterRequestBody
}

type RegisterRequestBody struct {
	Email       string `json:"email" required:"true" format:"email" example:"johndoe@example.com" minLength:"3"`
	Username    string `json:"username" required:"true" example:"johndoe" minLength:"3"`
	FirstName   string `json:"first_name" required:"true" example:"John" minLength:"1"`
	LastName    string `json:"last_name" required:"true" example:"Doe" minLength:"1"`
	PhoneNumber string `json:"phone_number" required:"false" minLength:"1"`
	Password    string `json:"password" required:"true" example:"password" minLength:"6"`
}

type RegisterResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
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
		switch {
		case errors.Is(err, domain.ErrUserDuplicateEmail):
			return nil, huma.Error409Conflict("email already taken")

		case errors.Is(err, domain.ErrUserDuplicateUsername):
			return nil, huma.Error409Conflict("username already taken")

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	respData := responses.NewUser.FromDomainUser(createdUser)

	otp, err := h.app.Services.Auth.GenerateOtp(ctx, createdUser.Id)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	if err := h.app.Services.User.SendVerificationEmail(ctx, createdUser, otp); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	return &RegisterResponse{
		Body: responses.Success(respData),
	}, nil
}
