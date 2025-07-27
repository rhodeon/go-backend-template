package contextutils

import (
	"context"
)

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, contextKeyRequestId, requestId)
}
