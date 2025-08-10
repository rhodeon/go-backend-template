package log

import (
	"fmt"
	"log/slog"

	"github.com/go-errors/errors"
)

// formatErrWithStackTrace returns a slog.GroupValue with keys `msg` and `trace`.
// If the error does not implement interface { StackTrace() errors.Error }, the `trace` key is omitted.
func formatErrWithStackTrace(err error) slog.Value {
	var groupedValues []slog.Attr

	// The error message is always set.
	groupedValues = append(groupedValues, slog.String(AttrErrorMsg, err.Error()))

	// A depth of 5 should be enough to cover the relevant traces in most cases and limits the noise from internal libraries.
	// It can be increased or removed entirely if it proves to be insufficient.
	const depth = 5

	var goErr *errors.Error
	if errors.As(err, &goErr) {
		var errorTrace []string

		for i, frame := range goErr.StackFrames() {
			if i >= depth {
				break
			}
			errorTrace = append(errorTrace, fmt.Sprintf("%s:%d", frame.File, frame.LineNumber))
		}

		groupedValues = append(groupedValues, slog.Any(AttrErrorTrace, errorTrace))
	}

	return slog.GroupValue(groupedValues...)
}
