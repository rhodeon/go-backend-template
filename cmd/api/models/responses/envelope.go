package responses

import (
	"strings"

	"github.com/rhodeon/go-backend-template/cmd/api/models/common"
)

// Envelope is a uniform wrapper around the actual response data that can hold more information like metadata and pagination.
type Envelope[T common.OasSchema] struct {
	Body ResponseBody[T]
}

func (e *Envelope[T]) Name() string {
	return strings.TrimSuffix(e.Body.Name(), "Body")
}

// Success wraps the passed in data into the uniform envelop form.
func Success[T common.OasSchema](data T) *Envelope[T] {
	return &Envelope[T]{
		Body: ResponseBody[T]{
			Data: data,
		},
	}
}

// ResponseBody is a shortcut to skip manually declaring the body for all responses.
type ResponseBody[T common.OasSchema] struct {
	Data T `json:"data"`
}

func (r ResponseBody[T]) Name() string {
	return strings.TrimSuffix(r.Data.Name(), "Data") + "Body"
}

type EnvelopeWithMetadata[T common.OasSchema, U Metadata] struct {
	Body ResponseBodyWithMetadata[T, U]
}

// ResponseBodyWithMetadata is an extension of ResponseBody with metadata information.
type ResponseBodyWithMetadata[T common.OasSchema, U Metadata] struct {
	Data     T `json:"data"`
	Metadata U `json:"metadata"`
}

func (r ResponseBodyWithMetadata[T, U]) Name() string {
	return strings.TrimSuffix(r.Data.Name(), "Data") + "Body"
}

func SuccessWithMetadata[T common.OasSchema, U Metadata](data T, metadata U) *EnvelopeWithMetadata[T, U] {
	re := &EnvelopeWithMetadata[T, U]{
		ResponseBodyWithMetadata[T, U]{
			Data:     data,
			Metadata: metadata,
		},
	}

	return re
}
