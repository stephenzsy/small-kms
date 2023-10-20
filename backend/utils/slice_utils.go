package utils

type MapFunc[T any, U any] func(U) T

func MapSlice[T any, U any](from []U, mapFunc MapFunc[T, U]) []T {
	if from == nil {
		return nil
	}
	to := make([]T, len(from))
	for i, v := range from {
		to[i] = mapFunc(v)
	}
	return to
}

func FilterSlice[T any](from []T, f func(T) bool) []T {
	if from == nil {
		return nil
	}
	to := make([]T, 0, len(from))
	for _, v := range from {
		if f(v) {
			to = append(to, v)
		}
	}
	return to
}

func NilIfZeroLen[T any](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	return s
}

func ToMapFunc[T any, K comparable](s []T, keyFunc func(item T) K) map[K]T {
	m := make(map[K]T, len(s))
	for _, item := range s {
		m[keyFunc(item)] = item
	}
	return m
}

func ToValueMapFunc[T any, K comparable, U any](s []T, keyValueFunc func(item T) (K, U)) map[K]U {
	m := make(map[K]U, len(s))
	for _, item := range s {
		k, v := keyValueFunc(item)
		m[k] = v
	}
	return m
}
