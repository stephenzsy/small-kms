package utils

func MapSlices[T any, U any](from []U, f func(U) T) []T {
	if from == nil {
		return nil
	}
	to := make([]T, len(from))
	for i, v := range from {
		to[i] = f(v)
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
