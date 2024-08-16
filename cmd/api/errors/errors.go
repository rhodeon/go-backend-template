package errors

import (
	"github.com/danielgtaylor/huma/v2"
	"log/slog"
)

func InternalServerError(err error, logger *slog.Logger) error {
	logger.Error("internal server error", slog.Any("error", err.Error()))
	return huma.Error500InternalServerError("an unexpected error occurred")
}
