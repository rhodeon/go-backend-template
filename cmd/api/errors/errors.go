package errors

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
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
			logMessage := "internal server error"
			if len(errs) > 0 {
				// If the error was due to a panic, the log message changes to reflect that.
				if len(errs) == 2 && errors.Is(errs[1], ErrPanic) {
					logMessage = "panic"
				}
				logger.Error(logMessage, slog.Any("error", errs[0].Error()))
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
