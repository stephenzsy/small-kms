package shared

import (
	"encoding"
	"encoding/base64"
	"encoding/hex"
	"net"
	"slices"
)

type Base64RawURLEncodableBytes []byte

// MarshalText implements encoding.TextMarshaler.
func (b Base64RawURLEncodableBytes) MarshalText() (text []byte, err error) {
	text = make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(text, b)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b Base64RawURLEncodableBytes) UnmarshalText(text []byte) error {
	b = make([]byte, base64.RawURLEncoding.DecodedLen(len(text)))
	_, err := base64.RawURLEncoding.Decode(b, text)
	return err
}

var _ encoding.TextMarshaler = Base64RawURLEncodableBytes{}
var _ encoding.TextUnmarshaler = (Base64RawURLEncodableBytes)(nil)

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

var _ encoding.TextMarshaler = net.IP{}
var _ encoding.TextUnmarshaler = (*net.IP)(nil)
