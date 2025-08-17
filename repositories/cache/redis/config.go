package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host        string
	Port        int
	Password    string
	Database    int
	OtpDuration time.Duration
	OnConnect   func(ctx context.Context, cn *redis.Conn) error
}
