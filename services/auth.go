package services

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/rhodeon/go-backend-template/repositories"

	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
)

type authTokenType string

const (
	AuthTokenTypeAccess  authTokenType = "access"
	AuthTokenTypeRefresh authTokenType = "refresh"
)

type Auth struct {
	*service
}

var authService *Auth

func newAuth(repos *repositories.Repositories, cfg *Config) *Auth {
	authService = &Auth{
		newService(repos, cfg),
	}
	return authService
}

// GenerateOtp generates a 6-digit one-time password.
func (a *Auth) GenerateOtp(ctx context.Context, userId int64) (string, error) {
	maxValue := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return "", errors.Errorf("generating otp: %w", err)
	}
	code := n.Int64() + 100000
	codeStr := strconv.FormatInt(code, 10)

	if err := a.repos.Cache.SetOtp(ctx, codeStr, userId); err != nil {
		return "", errors.Errorf("caching otp: %w", err)
	}

	return strconv.FormatInt(code, 10), nil
}

// ValidateOtp verifies if the provided OTP matches the stored code for the given user ID and clears it upon success.
func (a *Auth) ValidateOtp(ctx context.Context, code string, userId int64) (bool, error) {
	savedId, exists, err := a.repos.Cache.GetUserIdFromOtp(ctx, code)
	if err != nil {
		return false, errors.Errorf("retrieving user id from otp: %w", err)
	}
	if !exists {
		return false, nil
	}

	valid := savedId == userId
	if valid {
		if err := a.repos.Cache.ClearOtp(ctx, code); err != nil {
			return false, errors.Errorf("clearing otp: %w", err)
		}
	}

	return valid, nil
}

func (a *Auth) hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.Errorf("generating password with bcrypt: %w", err)
	}

	return string(passwordHash), nil
}
