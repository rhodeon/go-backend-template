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
			attr.Value = fmtErr(v)

		default:
			// Does nothing but satisfy linter.
		}

	default:
		// Does nothing but satisfy linter.
	}

	return attr
}
