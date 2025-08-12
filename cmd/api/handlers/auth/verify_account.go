package auth

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-errors/errors"
)

type VerifyAccountRequest struct {
	Body VerifyAccountRequestBody
}

type VerifyAccountRequestBody struct {
	Email string `json:"email" required:"true" format:"email" example:"johndoe@example.com" minLength:"3"`
	Otp   string `json:"otp" required:"true" example:"123456" minlength:"6" maxlength:"6"`
}

type VerifyAccountResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *Handlers) verifyAccount(ctx context.Context, req *VerifyAccountRequest) (*VerifyAccountResponse, error) {
	dbTx, commit, rollback, err := h.app.Db.BeginTx(ctx)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	if err := h.app.Services.User.Verify(ctx, dbTx, req.Body.Email, req.Body.Otp); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return nil, huma.Error404NotFound("User not found")

		case errors.Is(err, domain.ErrUserAlreadyVerified):
			return nil, huma.Error409Conflict("User is already verified")

		case errors.Is(err, domain.ErrAuthInvalidOtp):
			return nil, huma.Error401Unauthorized("Invalid OTP")

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	return &VerifyAccountResponse{
		Body: responses.Success[responses.SuccessMessage]("Success"),
	}, nil
}
