package errors

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/utils/contextutils"

	"github.com/danielgtaylor/huma/v2"
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

// NewApiError is a reimplementation of the default huma.NewErrorWithContext and allows ApiError to be used as the default over huma.ErrorModel.
func NewApiError() func(int, string, ...error) huma.StatusError {
	return func(status int, msg string, errs ...error) huma.StatusError {
		// Internal server errors have their underlying error changed to a generic message as
		// clients do not need to know the inner workings of the API.
		if status == http.StatusInternalServerError {
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

// UntypedError should be called in the default case of errors in the handlers after all recognised have been handled.
// Such an error could be a generic internal server error, a cancelled session, or could be due to the request timeout being exceeded.
// A 504 Gateway Timeout (server timeout) is returned in the latter case.
func UntypedError(ctx context.Context, err error) error {
	logger := contextutils.GetLogger(ctx)

	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		// If the request times out, a warning is logged. A high volume of such warning should be worth investigating.
		logger.Warn("Server timed out")
		return huma.Error504GatewayTimeout("server timeout")

	case errors.Is(err, context.Canceled):
		// If the session is cancelled (either explicitly by the user or something else).
		// a warning is logged. A high volume of such warnings should be worth investigating.
		logger.Warn("Session cancelled")
		return huma.Error500InternalServerError("", err)

	default:
		logger.Error("Internal server error", slog.Any(log.AttrError, err))
		return huma.Error500InternalServerError("", err)
	}
}
