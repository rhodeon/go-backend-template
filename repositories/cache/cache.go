package cache

import "context"

type Cache interface {
	SetOtp(ctx context.Context, code string, userId int64) error
	GetUserIdFromOtp(ctx context.Context, code string) (int64, bool, error)
	ClearOtp(ctx context.Context, code string) error
}
