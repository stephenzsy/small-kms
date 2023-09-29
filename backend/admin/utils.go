package admin

import (
	"encoding/base64"
	"encoding/hex"
)

// Ptr returns a pointer to the provided value.
func ToPtr[T any](v T) *T {
	return &v
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
