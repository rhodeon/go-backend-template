package log

import (
	"github.com/rhodeon/go-backend-template/internal/config"
	"log/slog"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	if cfg.DebugMode {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	logger := slog.Default()
	return logger
}
