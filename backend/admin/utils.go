package admin

// Ptr returns a pointer to the provided value.
func ToPtr[T any](v T) *T {
	return &v
}
