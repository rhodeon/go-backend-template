package errors

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
	"github.com/rhodeon/go-backend-template/internal/log"
	"log/slog"
	"net/http"
)

// ApiError is a subset of the default huma.ErrorModel with both the status and title stripped out
// from the final response as they can be inferred from the header and are redundant.
type ApiError struct {
	status int
	Detail string              `json:"detail,omitempty" example:"username is taken" doc:"A human-readable explanation specific to this occurrence of the problem."`
	Errors []*huma.ErrorDetail `json:"errors,omitempty" doc:"Optional list of individual error details."`
}

func (e *ApiError) Error() string {
	return e.Detail
}

func (e *ApiError) GetStatus() int {
	return e.status
}

// NewApiError is a reimplementation of the default huma.NewError and allows ApiError to be used as the default over huma.ErrorModel.
func NewApiError(logger *slog.Logger) func(int, string, ...error) huma.StatusError {
	return func(status int, msg string, errs ...error) huma.StatusError {
		// Internal server errors have their underlying error changed to a generic message as
		// clients do not need to know the inner workings of the API.
		if status == http.StatusInternalServerError {
			if len(errs) > 0 {
				switch {
				case len(errs) == 2 && errors.Is(errs[1], ErrPanic):
					// If the error was due to a panic, the log message reflects that.
					logger.Error("panic", slog.Any(log.AttrError, errs[0]))

				case errors.Is(errs[0], context.Canceled):
					// If the session is cancelled (either explicitly by the user or something else).
					// a warning is logged. A high volume of such warning should be worth investigating.
					logger.Warn("session cancelled")

				default:
					logger.Error("internal server error", slog.Any(log.AttrError, errs[0]))
				}
			}

			return &ApiError{
				status: status,
				Detail: "internal server error",
			}
		}

		details := make([]*huma.ErrorDetail, len(errs))
		for i := 0; i < len(errs); i++ {
			if converted, ok := errs[i].(huma.ErrorDetailer); ok {
				details[i] = converted.ErrorDetail()
			} else {
				if errs[i] == nil {
					continue
				}
				details[i] = &huma.ErrorDetail{Message: errs[i].Error()}
			}
		}

		return &ApiError{
			status: status,
			Detail: msg,
			Errors: details,
		}
	}
}

// ErrPanic is used to signify that the current error being handled by NewApiError()
// is as a result of a panic.
var ErrPanic = errors.New("panic")

// HandleUntypedError should be called in the default case of errors in the handlers after all recognised have been handled.
// Such an error could be a generic internal server error or could be due to the request timeout being exceeded.
// A 504 Gateway Timeout (server timeout) is returned in the latter case.
func HandleUntypedError(ctx context.Context, err error) error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return huma.Error504GatewayTimeout("server timeout")
	}
	return huma.Error500InternalServerError("", err)
}
