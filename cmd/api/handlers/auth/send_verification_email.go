package auth

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-errors/errors"
)

type SendVerificationEmailRequest struct {
	Body SendVerificationEmailRequestBody
}

type SendVerificationEmailRequestBody struct {
	Email string `json:"email" required:"true" format:"email" example:"johndoe@example.com" minLength:"3"`
}

type SendVerificationEmailResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *Handlers) sendVerificationEmail(ctx context.Context, req *SendVerificationEmailRequest) (*SendVerificationEmailResponse, error) {
	dbTx, commit, rollback, err := h.app.Db.BeginTx(ctx)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	user, err := h.app.Services.User.GetByEmail(ctx, dbTx, req.Body.Email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return nil, huma.Error401Unauthorized("Unauthenticated")

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	if user.IsVerified {
		return nil, huma.Error409Conflict("User is already verified")
	}

	otp, err := h.app.Services.Auth.GenerateOtp(ctx, user.Id)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	if err := h.app.Services.User.SendVerificationEmail(ctx, user, otp); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	return &SendVerificationEmailResponse{
		Body: responses.Success[responses.SuccessMessage]("Success"),
	}, nil
}
