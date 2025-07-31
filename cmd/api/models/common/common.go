package common

// OasSchema allows structs to have custom schema names in the OpenAPI documentation.
type OasSchema interface {
	Name() string
}
