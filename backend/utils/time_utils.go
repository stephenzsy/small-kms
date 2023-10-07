package utils

import (
	"time"
)

func IsTimeNotNilOrZero(t *time.Time) bool {
	return t != nil && !t.IsZero()
}
