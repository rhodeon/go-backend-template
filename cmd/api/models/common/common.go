package common

// InlinedSchema implementations are forced to be inlined under their endpoints in the OpenAPI docs,
// rather than being stored as reusable components.
type InlinedSchema interface {
	InlinedSchema()
}
