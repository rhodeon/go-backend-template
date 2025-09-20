package log

import (
	"context"
	"log/slog"
	"os"

	slogctx "github.com/veqryn/slog-context"
	slogotel "github.com/veqryn/slog-context/otel"
)

// Setup builds a new logger which automatically attaches OTel trace and span ids to logs.
func Setup(debugMode bool) {
	slogotel.DefaultKeyTraceID = "trace_id"
	slogotel.DefaultKeySpanID = "span_id"

	logLevel := slog.LevelInfo
	if debugMode {
		logLevel = slog.LevelDebug
	}

	baseOptions := &slog.HandlerOptions{
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
	}
	baseHandler := slog.NewJSONHandler(os.Stdout, baseOptions)

	ctxHandler := slogctx.NewHandler(baseHandler, &slogctx.HandlerOptions{
		Appenders: []slogctx.AttrExtractor{
			slogotel.ExtractTraceSpanID,
			slogctx.ExtractAppended,
		},
	})

	logger := slog.New(ctxHandler)
	slog.SetDefault(logger)
}

func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	switch attr.Value.Kind() {
	case slog.KindAny:
		switch v := attr.Value.Any().(type) {
		case error:
			// Errors are formatted with a stacktrace.
			attr.Value = formatErrWithStackTrace(v)

		default:
			// Does nothing but satisfy linter.
		}

	default:
		// Does nothing but satisfy linter.
	}

	return attr
}

// Fatal is a convenience wrapper to log fatal errors and immediately exit the program with an error.
func Fatal(ctx context.Context, message string, attrs ...slog.Attr) {
	anyAttrs := []any{}
	for _, attr := range attrs {
		anyAttrs = append(anyAttrs, attr)
	}

	slog.ErrorContext(ctx, message, anyAttrs...) //nolint:sloglint
	os.Exit(1)
}
