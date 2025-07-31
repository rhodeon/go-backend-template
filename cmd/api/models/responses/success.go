package responses

import "github.com/danielgtaylor/huma/v2"

// SuccessMessage is a wrapper to allow standalone string responses be documented.
type SuccessMessage string

func (rd SuccessMessage) Schema(_ huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:        "string",
		Nullable:    false,
		Description: "Standalone success message.",
		Examples:    []any{"Success"},
	}
}
