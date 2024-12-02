package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strings"
)

// tracer implements the pgx.QueryTracer interface to provider debugging and tracing capabilities for queries.
type tracer struct {
	logger    *slog.Logger
	debugMode bool
}

func newTracer(logger *slog.Logger, debugMode bool) tracer {
	return tracer{
		logger:    logger,
		debugMode: debugMode,
	}
}

// TraceQueryStart logs each database query triggered when in debug mode.
// The debugMode is used rather than simply calling the logger.Debug in order to
// prevent wasting time on formatting the output when not in debug mode.
func (t tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if t.debugMode {
		// Prevent pollution of logs with transaction boundary commands.
		if strings.EqualFold(data.SQL, "begin") ||
			strings.EqualFold(data.SQL, "rollback") ||
			strings.EqualFold(data.SQL, "commit") ||
			strings.EqualFold(data.SQL, "end") {
			return ctx
		}
		var render = "\n"
		for i, arg := range data.Args {
			switch arg.(type) {
			case []byte:
				render += fmt.Sprintf("$%d:\t%s\n", i+1, arg)
			default:
				render += fmt.Sprintf("$%d:\t%v\n", i+1, arg)
			}
		}

		t.logger.Info(fmt.Sprintf("Executing db query:\n%s\nargs:%s", data.SQL, render))
	}

	return ctx
}

func (t tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {}
