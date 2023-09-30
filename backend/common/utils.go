package common

import "github.com/google/uuid"

/*
type resolvePtrWithDefault[D any] struct {
	ptr          *D
	defaultValue D
}

func (w *resolvePtrWithDefault[D]) Value() D {
	if w.ptr == nil {
		return w.defaultValue
	}
	return *w.ptr
}

type ResolvePtr[D any] interface {
	Value() D
}

func ResolvePtrWithDefault[D any](ptr *D, defaultValue D) ResolvePtr[D] {
	return &resolvePtrWithDefault[D]{ptr: ptr, defaultValue: defaultValue}
}
*/
func ResolveBoolPtrValue(ptr *bool) bool {
	return ptr != nil && *ptr
}

func UUIDWithinRange(id, start, end uuid.UUID) bool {
	for i, b := range id {
		if b < start[i] || b > end[i] {
			return false
		}
	}
	return true
}
