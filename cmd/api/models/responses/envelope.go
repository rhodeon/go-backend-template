package responses

// Envelope is a uniform wrapper around the actual response data that can hold more information like metadata and pagination.
type Envelope[T any] struct {
	Body ResponseBody[T]
}

// Success wraps the passed in data into the uniform envelop form.
func Success[T any](data T) *Envelope[T] {
	return &Envelope[T]{
		Body: ResponseBody[T]{
			Data: data,
		},
	}
}

// ResponseBody is a shortcut to skip manually declaring the body for all responses.
type ResponseBody[T any] struct {
	Data T `json:"data"`
}

type EnvelopeWithMetadata[T any, U any] struct {
	Body ResponseBodyWithMetadata[T, U]
}

// ResponseBodyWithMetadata is an extension of ResponseBody with metadata information.
type ResponseBodyWithMetadata[T any, U any] struct {
	Data     T `json:"data"`
	Metadata U `json:"metadata"`
}

func SuccessWithMetadata[T any, U any](data T, metadata U) *EnvelopeWithMetadata[T, U] {
	return &EnvelopeWithMetadata[T, U]{
		Body: ResponseBodyWithMetadata[T, U]{
			Data:     data,
			Metadata: metadata,
		},
	}
}
