package middleware

import (
	"strings"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/utils/contextutils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// RequestId adds a unique ID to associate with user issues.
func RequestId(_ *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		requestId := strings.ReplaceAll(uuid.New().String(), "-", "")
		newCtx := contextutils.WithRequestId(ctx.Context(), requestId)

		ctx = huma.WithContext(ctx, newCtx)
		next(ctx)
	}
}
