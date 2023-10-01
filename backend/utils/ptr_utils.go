package utils

func ToPtr[D any](v D) *D {
	return &v
}

func NilToDefault[D any](ptr *D) (v D) {
	if ptr != nil {
		v = *ptr
	}
	return
}
