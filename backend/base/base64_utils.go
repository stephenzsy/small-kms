package base

import (
	"encoding"
	"encoding/base64"
)

type base64RawURLEncodedBytesImpl []byte

// MarshalText implements encoding.TextMarshaler.
func (b base64RawURLEncodedBytesImpl) MarshalText() (text []byte, err error) {
	text = make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(text, b)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b Base64RawURLEncodedBytes) UnmarshalText(text []byte) error {
	b = make([]byte, base64.RawURLEncoding.DecodedLen(len(text)))
	_, err := base64.RawURLEncoding.Decode(b, text)
	return err
}

var _ encoding.TextMarshaler = base64RawURLEncodedBytesImpl{}
var _ encoding.TextUnmarshaler = (base64RawURLEncodedBytesImpl)(nil)
