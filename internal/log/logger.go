package log

import (
	"log/slog"
)

func NewLogger(debugMode bool) *slog.Logger {
	if debugMode {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	logger := slog.Default()
	return logger
}
