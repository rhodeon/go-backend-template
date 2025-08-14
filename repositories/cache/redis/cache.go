package redis

import (
	"context"
	"strconv"

	"github.com/rhodeon/go-backend-template/repositories/cache"

	"github.com/go-errors/errors"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	config *Config
}

func New(ctx context.Context, cfg *Config) (cache.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Errorf("pinging redis: %w", err)
	}

	return &Cache{
		client: client,
		config: cfg,
	}, nil
}

func (c *Cache) SetOtp(ctx context.Context, code string, userId int64) error {
	if err := c.client.Set(ctx, c.buildOtpKey(code), userId, c.config.OtpDuration).Err(); err != nil {
		return errors.Errorf("setting otp in redis: %w", err)
	}

	return nil
}

func (c *Cache) GetUserIdFromOtp(ctx context.Context, code string) (int64, bool, error) {
	userId, err := c.client.Get(ctx, c.buildOtpKey(code)).Int64()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return 0, false, nil

		default:
			return 0, false, errors.Errorf("getting user from otp in redis: %w", err)
		}
	}

	return userId, true, nil
}

func (c *Cache) ClearOtp(ctx context.Context, code string) error {
	if err := c.client.Del(ctx, c.buildOtpKey(code)).Err(); err != nil {
		return errors.Errorf("deleting otp from redis: %w", err)
	}

	return nil
}

func (c *Cache) buildOtpKey(code string) string {
	return namespaceOtp + ":" + code
}
