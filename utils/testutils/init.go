package testutils

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/utils/contextutils"
)

var projectRootDir string

func init() {
	logger := contextutils.GetLogger(context.Background())
	var err error
	if projectRootDir, err = getProjectRootDir(); err != nil {
		log.Fatal(logger, "Failed to get project root directory", slog.Any(log.AttrError, err))
	}

	if config, err = parseConfig(); err != nil {
		log.Fatal(logger, "Failed to parse config", slog.Any(log.AttrError, err))
	}
}
