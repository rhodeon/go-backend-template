package internal

import "log/slog"

func SetupLogger(cfg *Config) *slog.Logger {
	if cfg.DebugMode {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	logger := slog.Default()
	return logger
}
