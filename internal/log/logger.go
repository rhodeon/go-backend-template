package log

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
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
		// Errors are formatted with a stacktrace.
		switch v := attr.Value.Any().(type) {
		case error:
			attr.Value = fmtErr(v)

		default:
			// Does nothing but satisfy linter.
		}

	default:
		// Does nothing but satisfy linter.
	}

	return attr
}
