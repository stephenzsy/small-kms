package utils

type Nullable[T interface{}] struct {
	hasValue bool
	value    T
}

func (n *Nullable[T]) HasValue() bool {
	return n.hasValue
}

func (n *Nullable[T]) Value() T {
	return n.value
}

func (n *Nullable[T]) SetValue(value T) {
	n.hasValue = true
	n.value = value
}
