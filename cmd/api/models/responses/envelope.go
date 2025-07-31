package responses

// Envelope is a wrapper around the actual response data that can be extended to hold metadata like pagination details.
type Envelope[T any] struct {
	Data T `json:"data"`
}

// InlinedSchema ensures all response bodies are inlined in their specific endpoint paths in the OpenAPI docs.
func (b Envelope[any]) InlinedSchema() {}

// Success is a shortcut to skip manually declaring the body for all successful responses.
func Success[T any](data T) Envelope[T] {
	return Envelope[T]{
		Data: data,
	}
}

// EnvelopeWithMetadata is an extension of Envelope with metadata information.
type EnvelopeWithMetadata[T any, U any] struct {
	Envelope[T]
	Metadata U `json:"metadata"`
}

func SuccessWithMetadata[T any, U any](data T, metadata U) EnvelopeWithMetadata[T, U] {
	return EnvelopeWithMetadata[T, U]{
		Success(data),
		metadata,
	}
}
