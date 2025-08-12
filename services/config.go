package services

import "time"

type Config struct {
	Auth *AuthConfig
}

type AuthConfig struct {
	JwtIssuer               string
	JwtAccessTokenSecret    string
	JwtRefreshTokenSecret   string
	JwtAccessTokenDuration  time.Duration
	JwtRefreshTokenDuration time.Duration
	OtpDuration             time.Duration
}
