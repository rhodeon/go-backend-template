package cache

import "context"

type Cache interface {
	SetOtp(ctx context.Context, userId int64, code string) error
	GetOtp(ctx context.Context, userId int64) (string, bool, error)
	ClearOtp(ctx context.Context, userId int64) error
}
