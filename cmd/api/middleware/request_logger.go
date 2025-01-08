package middleware

import (
	"log/slog"
	"time"

	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rhodeon/go-backend-template/internal/helpers"
)

// RequestLogger logs details associated with the request, like the HTTP method, IP address, response status code, e.t.c.
func RequestLogger() func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		startTime := time.Now()

		defer func() {
			logger := helpers.ContextGetLogger(ctx.Context())
			logger.Info(
				"[REQ]",
				slog.String(log.AttrMethod, ctx.Method()),
				slog.String(log.AttrPath, ctx.URL().Path),
				slog.String(log.AttrIP, ctx.RemoteAddr()),
				slog.Int(log.AttrStatus, ctx.Status()),
				slog.Any(log.AttrDuration, time.Since(startTime)),
			)
		}()

		next(ctx)
	}
}
