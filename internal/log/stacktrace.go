package log

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// fmtErr returns a slog.GroupValue with keys "msg" and "trace". If the error
// does not implement interface { StackTrace() errors.StackTrace }, the "trace"
// key is omitted.
func fmtErr(err error) slog.Value {
	var groupedValues []slog.Attr

	groupedValues = append(groupedValues, slog.String(AttrErrorMsg, err.Error()))

	// The earliest underlying error with a stacktrace is determined by iteratively unwrapping each error until the cause is hit.
	// errors.Cause() is not used directly, because errors with stack traces would be missed if the underlying cause
	// does not have a trace attached.
	var st stackTracer
	for err := err; err != nil; err = errors.Unwrap(err) {
		if x, ok := err.(stackTracer); ok {
			st = x
		}
	}

	if st != nil {
		traceLines := getStacktraceLines(st.StackTrace())
		groupedValues = append(groupedValues, slog.Any(AttrErrorTrace, traceLines))
	}

	return slog.GroupValue(groupedValues...)
}

func getStacktraceLines(frames errors.StackTrace) []string {
	traceLines := make([]string, len(frames))

	// Iterate in reverse to skip unnecessary, consecutive runtime frames at the bottom of the trace.
	var skipped int
	skipping := true
	for i := len(frames) - 1; i >= 0; i-- {
		// Adapted from errors.Frame.MarshalText(), but avoiding repeated calls to FuncForPC and FileLine.
		pc := uintptr(frames[i]) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			traceLines[i] = "unknown"
			skipping = false
			continue
		}

		functionName := fn.Name()

		if skipping && strings.HasPrefix(functionName, "runtime.") {
			skipped++
			continue
		} else {
			skipping = false
		}

		filename, lineNr := fn.FileLine(pc)
		traceLines[i] = fmt.Sprintf("'%s -> %s:%d'", functionName, filename, lineNr)
	}

	return traceLines[:len(traceLines)-skipped]
}
