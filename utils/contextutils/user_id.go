package contextutils

import "context"

func WithUserId(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, contextKeyUserId, userId)
}
