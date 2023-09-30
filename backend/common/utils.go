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

func SliceMap[T any, U any](slice []T, mapper func(T) U) []U {
	if slice == nil {
		return nil
	}
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)

	}
	return result
}

func SliceMapWithError[T any, U any](slice []T, mapper func(T) (U, error)) ([]U, error) {
	if slice == nil {
		return nil, nil
	}
	result := make([]U, len(slice))
	var err error
	for i, item := range slice {
		result[i], err = mapper(item)
		if err != nil {
			return result, err
		}

	}
	return result, err
}
