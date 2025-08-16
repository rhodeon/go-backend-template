package temp

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
)

var ProjectRootDir string

// init sets up the resources needed before running integration tests.
func init() {
	var err error
	if ProjectRootDir, err = getProjectRootDir(); err != nil {
		log.Fatal(err)
	}
}

// getProjectRootDir determines the root directory of the project by finding the first location with the go.mod file.
func getProjectRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.Errorf("getting current working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found go.mod, this is the root.
			return dir, nil
		}

		// If the project root hasn't been found, the parent directory is checked next.
		if parent := filepath.Dir(dir); parent == dir {
			// The system root has been reached without any go.mod file found.
			return "", errors.New("go.mod file not found")
		} else {
			dir = parent
		}
	}
}
