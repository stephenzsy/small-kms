package admin

import (
	"encoding/base64"
	"encoding/hex"
)

// Ptr returns a pointer to the provided value.
func ToPtr[T any](v T) *T {
	return &v
}

func DefaultIfNil[D any](ptr *D) (value D) {
	if ptr != nil {
		value = *ptr
	}
	return
}

func ToOptionalStringPtr(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func base64UrlToHexStr(s string) string {
	if b, err := base64.URLEncoding.DecodeString(s); err == nil {
		return hex.EncodeToString(b)
	}
	return ""
}

func base64UrlToHexStrPtr(s *string) *string {
	str := base64UrlToHexStr(*s)
	if str == "" {
		return nil
	}
	return ToPtr(str)
}
