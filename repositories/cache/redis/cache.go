package redis

import (
	"context"
	"fmt"
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

func (c *Cache) SetOtp(ctx context.Context, userId int64, code string) error {
	if err := c.client.Set(ctx, c.buildOtpKey(userId), code, c.config.OtpDuration).Err(); err != nil {
		return errors.Errorf("setting otp in redis: %w", err)
	}

	return nil
}

func (c *Cache) GetOtp(ctx context.Context, userId int64) (string, bool, error) {
	otp, err := c.client.Get(ctx, c.buildOtpKey(userId)).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return "", false, nil

		default:
			return "", false, errors.Errorf("getting otp from in redis: %w", err)
		}
	}

	return otp, true, nil
}

func (c *Cache) ClearOtp(ctx context.Context, userId int64) error {
	if err := c.client.Del(ctx, c.buildOtpKey(userId)).Err(); err != nil {
		return errors.Errorf("deleting otp from redis: %w", err)
	}

	return nil
}

func (c *Cache) buildOtpKey(userId int64) string {
	return fmt.Sprintf("%s:%d", namespaceOtp, userId)
}
