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
	ErrUserAlreadyVerified   = errors.New("")

	ErrAuthExpiredToken = errors.New("auth token has expired")
	ErrAuthInvalidToken = errors.New("auth token is invalid")
	ErrAuthInvalidOtp   = errors.New("otp is invalid")
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

func UserErrAlreadyVerified(id int64) error {
	return embedBaseErrorf(ErrUserAlreadyVerified, "user with id %q is already verified", id)
}

func UserErrDuplicateUsername(username string) error {
	return embedBaseErrorf(ErrUserDuplicateUsername, "user with username %q already exists", username)
}
