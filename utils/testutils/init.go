package testutils

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/log"
)

var projectRootDir string

func init() {
	var err error
	if projectRootDir, err = getProjectRootDir(); err != nil {
		log.Fatal(context.Background(), "Failed to get project root directory", slog.Any(log.AttrError, err))
	}

	if config, err = parseConfig(); err != nil {
		log.Fatal(context.Background(), "Failed to parse config", slog.Any(log.AttrError, err))
	}
}
