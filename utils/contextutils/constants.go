package contextutils

type contextKey string

const (
	contextKeyLogger    contextKey = "logger"
	contextKeyRequestId contextKey = "request_id"
	contextKeyUserId    contextKey = "user_id"
)
