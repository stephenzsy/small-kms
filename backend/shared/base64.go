package shared

import (
	"encoding"
	"encoding/base64"
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
