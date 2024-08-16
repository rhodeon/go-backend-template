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

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.RFC3339,
		}),
	))

	logger := slog.Default()
	return logger
}
