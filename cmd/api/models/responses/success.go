package responses

import "github.com/danielgtaylor/huma/v2"

type SuccessMessageResponseData string

func (rd SuccessMessageResponseData) Name() string {
	return "SuccessMessageResponseData"
}

func (rd SuccessMessageResponseData) Schema(_ huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:        "string",
		Nullable:    false,
		Description: "Standalone success message",
		Examples:    []any{"Success"},
	}
}
