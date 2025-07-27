package typeutils

// Ptr returns a pointer to the given value. Useful for creating pointers to literals in a concise manner.
func Ptr[T any](value T) *T {
	return &value
}

// Deref returns the value pointed to by ptr, or the zero value of T if ptr is nil.
func Deref[T any](ptr *T) T {
	if ptr == nil {
		return *new(T)
	}

	return *ptr
}
