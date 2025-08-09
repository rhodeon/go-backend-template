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
		Level:      logLevel,
		TimeFormat: time.RFC3339,
	}
	tintHandler := tint.NewHandler(os.Stderr, tintOptions)

	logger := slog.New(tintHandler)
	slog.SetDefault(logger)
	return logger
}
