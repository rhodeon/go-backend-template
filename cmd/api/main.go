package main

import (
	"context"
	"log/slog"
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
	if err := run(mainCtx, cfg); err != nil {
		log.Fatal(mainCtx, "Error running API", slog.Any(log.AttrError, err))
	}
}

// run is separated from the main function to ensure cleanups are honoured when an error occurs.
// Calling log.Fatal directly in run would trigger os.Exit which skips deferred functions.
func run(ctx context.Context, cfg *internal.Config) error {
	log.Setup(cfg.DebugMode)

	otelConfig := otel.Config(cfg.Otel)
	if err := otel.SetupProviders(ctx, &otelConfig); err != nil {
		return errors.Errorf("Setting up OTel providers: %w", err)
	}

	dbConfig := database.Config(cfg.Database)
	db, closeDb, err := database.Connect(ctx, &dbConfig, cfg.DebugMode)
	if err != nil {
		return errors.Errorf("connecting to database: %w", err)
	}
	defer closeDb()

	repos, err := setupRepositories(ctx, cfg)
	if err != nil {
		return errors.Errorf("setting up repositories: %w", err)
	}
	svcs := setupServices(ctx, cfg, repos)
	app := internal.NewApplication(cfg, db, svcs)

	// A waitgroup is established to ensure background tasks are completed before shutting down the server.
	backgroundWg := &sync.WaitGroup{}

	// The listen chan isn't used here and is buffered to 1 so the server won't be blocked.
	if err = server.ServeApi(ctx, app, backgroundWg, make(chan<- int, 1)); err != nil {
		return errors.Errorf("serving API: %w", err)
	}

	return nil
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
	// 	Host:            cfg.Smtp.Host,
	// 	Port:            cfg.Smtp.Port,
	// 	User:            cfg.Smtp.User,
	// 	Password:        cfg.Smtp.Password,
	// 	Sender:          cfg.Smtp.Sender,
	// 	OtpDuration:     cfg.Auth.OtpDuration,
	// 	OtelServiceName: cfg.Otel.ServiceName,
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
