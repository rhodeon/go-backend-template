package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/rhodeon/go-backend-template/cmd/api/internal"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
)

// Logger logs details associated with the request, like the HTTP method, IP address, response status code, etc.
func Logger(_ *internal.Application, _ huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		startTime := time.Now()

		defer func() {
			slog.InfoContext(
				ctx.Context(),
				fmt.Sprintf("[REQ] %s %s", ctx.Method(), ctx.URL().Path), //nolint:sloglint
				slog.String(log.AttrMethod, ctx.Method()),
				slog.String(log.AttrPath, ctx.URL().Path),
				slog.String(log.AttrIP, ctx.RemoteAddr()),
				slog.Int(log.AttrStatus, ctx.Status()),
				slog.String(log.AttrDuration, time.Since(startTime).String()),
			)
		}()

		next(ctx)
	}
}
