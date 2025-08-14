package auth

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Body LoginRequestBody
}

type LoginRequestBody struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"123456"`
}

type LoginResponse struct {
	Body responses.Envelope[LoginResponseData]
}

type LoginResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handlers) login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	dbTx, commit, rollback, err := h.app.Db.BeginTx(ctx)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	user, err := h.app.Services.User.GetByUsername(ctx, dbTx, req.Body.Username)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return nil, huma.Error401Unauthorized("Unauthenticated")

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	if !user.IsVerified {
		return nil, huma.Error401Unauthorized("User is unverified")
	}

	if err := h.app.Services.User.VerifyPassword(req.Body.Password, user.Password); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, huma.Error401Unauthorized("Unauthenticated")

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	accessToken, err := h.app.Services.Auth.GenerateAccessToken(user.Id)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	refreshToken, err := h.app.Services.Auth.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	respData := LoginResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &LoginResponse{
		Body: responses.Success(respData),
	}, nil
}
