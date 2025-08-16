package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func NewLogger(debugMode bool) *slog.Logger {
	logLevel := slog.LevelInfo
	if debugMode {
		logLevel = slog.LevelDebug
	}

	tintOptions := &tint.Options{
		Level:       logLevel,
		TimeFormat:  time.RFC3339,
		ReplaceAttr: replaceAttr,
	}
	tintHandler := tint.NewHandler(os.Stderr, tintOptions)

	_ = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
	})

	logger := slog.New(tintHandler)
	slog.SetDefault(logger)
	return logger
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
func Fatal(logger *slog.Logger, message string, attrs ...slog.Attr) {
	anyAttrs := []any{}
	for _, attr := range attrs {
		anyAttrs = append(anyAttrs, attr)
	}

	logger.Error(message, anyAttrs...) //nolint:sloglint
	os.Exit(1)
}
