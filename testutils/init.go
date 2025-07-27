package testutils

import (
	"context"
	"log"
)

var projectRootDir string

// init sets up the resources needed before running integration tests.
func init() {
	var err error
	if projectRootDir, err = getProjectRootDir(); err != nil {
		log.Fatal(err)
	}

	if config, err = parseConfig(); err != nil {
		log.Fatal(err)
	}

	if err := setupDatabaseContainer(context.Background()); err != nil {
		log.Fatal(err)
	}
}
