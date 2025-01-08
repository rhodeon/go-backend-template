package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/rhodeon/go-backend-template/internal/helpers"

	"github.com/jackc/pgx/v5"
)

// tracer implements the pgx.QueryTracer interface to provider debugging and tracing capabilities for queries.
type tracer struct {
	debugMode bool
}

func newTracer(debugMode bool) tracer {
	return tracer{
		debugMode: debugMode,
	}
}

// TraceQueryStart logs each database query triggered when in debug mode.
// The debugMode is used rather than simply calling the logger.Debug in order to
// prevent wasting time on formatting the output when not in debug mode.
func (t tracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if t.debugMode {
		// Prevent pollution of logs with transaction boundary commands.
		if strings.EqualFold(data.SQL, "begin") ||
			strings.EqualFold(data.SQL, "rollback") ||
			strings.EqualFold(data.SQL, "commit") ||
			strings.EqualFold(data.SQL, "end") {
			return ctx
		}
		render := "\n"
		for i, arg := range data.Args {
			switch arg.(type) {
			case []byte:
				render += fmt.Sprintf("$%d:\t%s\n", i+1, arg)
			default:
				render += fmt.Sprintf("$%d:\t%v\n", i+1, arg)
			}
		}

		helpers.ContextGetLogger(ctx).Info(fmt.Sprintf("Executing db query:\n%s\nargs:%s", data.SQL, render))
	}

	return ctx
}

func (t tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {}
