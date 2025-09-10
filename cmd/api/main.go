package main

import (
	"context"
	"sync"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/cmd/api/server"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/internal/otel"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/cache/redis"
	mockemail "github.com/rhodeon/go-backend-template/repositories/email/mock"
	"github.com/rhodeon/go-backend-template/services"

	"github.com/go-errors/errors"
)

func main() {
	mainCtx := context.Background()
	cfg := internal.ParseConfig()
	logger := log.NewLogger(cfg.DebugMode)

	otelConfig := otel.Config(cfg.Otel)
	if err := otel.NewTracer(&otelConfig); err != nil {
		panic(err)
	}

	dbConfig := database.Config(cfg.Database)
	db, closeDb, err := database.Connect(mainCtx, &dbConfig, cfg.DebugMode)
	if err != nil {
		panic(err)
	}
	defer closeDb()

	repos, err := setupRepositories(mainCtx, cfg)
	if err != nil {
		panic(err)
	}
	svcs := setupServices(mainCtx, cfg, repos)
	app := internal.NewApplication(cfg, logger, db, svcs)

	// A waitgroup is established to ensure background tasks are completed before shutting down the server.
	backgroundWg := &sync.WaitGroup{}

	// The listen chan isn't used here and is buffered to 1 so the server won't be blocked.
	err = server.ServeApi(app, backgroundWg, make(chan<- int, 1))
	if err != nil {
		panic(err)
	}
}

func setupRepositories(ctx context.Context, cfg *internal.Config) (*repositories.Repositories, error) {
	cache, err := redis.New(ctx, &redis.Config{
		Host:        cfg.Cache.Host,
		Port:        cfg.Cache.Port,
		Password:    cfg.Cache.Password,
		Database:    cfg.Cache.Database,
		OtpDuration: cfg.Auth.OtpDuration,
	})
	if err != nil {
		return nil, errors.Errorf("setting up redis: %w", err)
	}

	email := mockemail.New()
	// email, err := smtp.New(ctx, &smtp.Config{
	// 	Host:        cfg.Smtp.Host,
	// 	Port:        cfg.Smtp.Port,
	// 	User:        cfg.Smtp.User,
	// 	Password:    cfg.Smtp.Password,
	// 	Sender:      cfg.Smtp.Sender,
	// 	OtpDuration: cfg.Auth.OtpDuration,
	// })
	// if err != nil {
	// 	return nil, errors.Errorf("setting up smtp email repo: %w", err)
	// }

	repos := repositories.New(cache, email)
	return repos, nil
}

func setupServices(_ context.Context, cfg *internal.Config, repos *repositories.Repositories) *services.Services {
	servicesCfg := &services.Config{
		Auth: &services.AuthConfig{
			JwtIssuer:               cfg.Auth.JwtIssuer,
			JwtAccessTokenSecret:    cfg.Auth.JwtAccessTokenSecret,
			JwtRefreshTokenSecret:   cfg.Auth.JwtRefreshTokenSecret,
			JwtAccessTokenDuration:  cfg.Auth.JwtAccessTokenDuration,
			JwtRefreshTokenDuration: cfg.Auth.JwtRefreshTokenDuration,
			OtpDuration:             cfg.Auth.OtpDuration,
		},
	}

	svcs := services.New(repos, servicesCfg)
	return svcs
}
