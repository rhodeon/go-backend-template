package testutils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// getProjectRootDir determines the root directory of the project by finding the first location with the go.mod file.
func getProjectRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "unable to get current working directory")
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found go.mod, this is the root.
			return dir, nil
		}

		// If the project root hasn't been found, the parent directory is checked next.
		if parent := filepath.Dir(dir); parent == dir {
			// The system root has been reached without any go.mod file found.
			return "", fmt.Errorf("go.mod not found")
		} else {
			dir = parent
		}
	}
}

// JoinUrlPath is a thin wrapper around url.JoinPath removing the returned error from the signature for easier use in tests.
func JoinUrlPath(base string, paths ...string) string {
	result, err := url.JoinPath(base, paths...)
	if err != nil {
		panic(err)
	}

	return result
}
