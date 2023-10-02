package utils

import "github.com/google/uuid"

func ToPtr[D any](v D) *D {
	return &v
}

func NilToDefault[D any](ptr *D) (v D) {
	if ptr != nil {
		v = *ptr
	}
	return
}

func NonNilUUID(id *uuid.UUID) (uuid.UUID, bool) {
	if id == nil {
		return uuid.Nil, false
	}
	return *id, *id != uuid.Nil
}
