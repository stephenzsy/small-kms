package utils

import "github.com/google/uuid"

// test if uuid is in range [start, end], inclusive
func IsUUIDInRange(id, start, end uuid.UUID) bool {
	for i, b := range id {
		if b < start[i] || b > end[i] {
			return false
		}
	}
	return true
}

// in case uuid.Nil is altered
func IsUUIDNil(id uuid.UUID) bool {
	for _, b := range id {
		if b != 0 {
			return false
		}
	}
	return true
}
