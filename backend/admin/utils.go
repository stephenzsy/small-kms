package admin

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/stephenzsy/small-kms/backend/graph"
)

// Ptr returns a pointer to the provided value.
func ToPtr[T any](v T) *T {
	return &v
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

func OdataTypeToNSType(odataType graph.MsGraphOdataType) NamespaceTypeShortName {
	switch odataType {
	case graph.MsGraphOdataTypeDevice:
		return NSTypeDevice
	case graph.MsGraphOdataTypeGroup:
		return NSTypeGroup
	case graph.MsGraphOdataTypeUser:
		return NSTypeUser
	case graph.MsGraphOdataTypeApplication:
		return NSTypeApplication
	case graph.MsGraphOdataTypeServicePrincipal:
		return NSTypeServicePrincipal
	}
	return NSTypeAny
}
