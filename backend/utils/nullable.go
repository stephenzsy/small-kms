package utils

type Nullable[T interface{}] struct {
	HasValue bool
	Value    T
}
