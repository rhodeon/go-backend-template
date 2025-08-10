package domain

import (
	"fmt"

	"github.com/go-errors/errors"
)

// Domain errors are represented as sentinels with empty strings as their builder functions are used to
// generate the actual error message with dynamic data for more context.
var (
	ErrUserNotFound          = errors.New("")
	ErrUserDuplicateEmail    = errors.New("")
	ErrUserDuplicateUsername = errors.New("")
)

// embedBaseErrorf injects the underlying domain error into a new error generated from the format and args.
// This embedding is needed for errors.Is() calls to behave correctly.
func embedBaseErrorf(baseErr error, format string, args ...any) error {
	allArgs := []any{baseErr}
	allArgs = append(allArgs, args...)
	return errors.Wrap(fmt.Errorf("%w"+format, allArgs...), 2) //nolint:forbidigo
}

func UserErrNotFound(field string, value any) error {
	return embedBaseErrorf(ErrUserNotFound, "user with %s %q not found", field, value)
}

func UserErrDuplicateEmail(email string) error {
	return embedBaseErrorf(ErrUserDuplicateEmail, "user with email %q already exists", email)
}

func UserErrDuplicateUsername(username string) error {
	return embedBaseErrorf(ErrUserDuplicateUsername, "user with username %q already exists", username)
}
