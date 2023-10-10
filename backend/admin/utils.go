package admin

// Ptr returns a pointer to the provided value.
func ToPtr[T any](v T) *T {
	return &v
}

func ToOptionalStringPtr(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}
