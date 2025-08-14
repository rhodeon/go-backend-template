package services

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/rhodeon/go-backend-template/domain"
	"github.com/rhodeon/go-backend-template/repositories"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
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

type JWTClaims struct {
	UserId int64  `json:"sub"`
	Type   string `json:"typ"`
	jwt.RegisteredClaims
}

func (a *Auth) GenerateAccessToken(userId int64) (string, error) {
	claims := JWTClaims{
		userId,
		string(AuthTokenTypeAccess),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.cfg.Auth.JwtAccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.cfg.Auth.JwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.cfg.Auth.JwtAccessTokenSecret))
	if err != nil {
		return "", errors.Errorf("signing access token: %w", err)
	}
	return signedToken, nil
}

func (a *Auth) GenerateRefreshToken(userId int64) (string, error) {
	claims := JWTClaims{
		userId,
		string(AuthTokenTypeRefresh),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.cfg.Auth.JwtRefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.cfg.Auth.JwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.cfg.Auth.JwtRefreshTokenSecret))
	if err != nil {
		return "", errors.Errorf("signing refresh token: %w", err)
	}
	return signedToken, nil
}

func (a *Auth) ParseToken(tokenString string, tokenType authTokenType) (*JWTClaims, error) {
	var secret string
	switch tokenType {
	case AuthTokenTypeAccess:
		secret = a.cfg.Auth.JwtAccessTokenSecret

	case AuthTokenTypeRefresh:
		secret = a.cfg.Auth.JwtRefreshTokenSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, JWTClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "expired"):
			return nil, domain.ErrAuthExpiredToken
		default:
			return nil, domain.ErrAuthInvalidToken
		}
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, domain.ErrAuthInvalidToken
	}
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
