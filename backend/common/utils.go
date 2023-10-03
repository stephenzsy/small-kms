package common

func ResolveBoolPtrValue(ptr *bool) bool {
	return ptr != nil && *ptr
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
