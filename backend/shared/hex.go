package shared

import (
	"encoding"
	"encoding/hex"
	"slices"
)

type certificateFingerprintImpl []byte

// MarshalText implements encoding.TextMarshaler.
func (s certificateFingerprintImpl) MarshalText() (text []byte, _ error) {
	text = make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(text, s)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *certificateFingerprintImpl) UnmarshalText(text []byte) (err error) {
	text = slices.DeleteFunc(text, func(c byte) bool {
		return c == ':'
	})
	*s = make([]byte, hex.DecodedLen(len(text)))
	_, err = hex.Decode(*s, text)
	return
}

func (s *certificateFingerprintImpl) HexString() string {
	if s == nil {
		return ""
	}
	return hex.EncodeToString(*s)
}

var _ encoding.TextMarshaler = certificateFingerprintImpl{}
var _ encoding.TextUnmarshaler = (*certificateFingerprintImpl)(nil)
