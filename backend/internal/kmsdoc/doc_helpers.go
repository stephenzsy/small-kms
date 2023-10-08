package kmsdoc

import (
	"encoding"
	"encoding/base64"
	"encoding/hex"
	"time"
)

type TimeStorable time.Time

// MarshalText implements encoding.TextMarshaler.
func (t TimeStorable) MarshalText() (text []byte, _ error) {
	return time.Time(t).UTC().MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *TimeStorable) UnmarshalText(text []byte) error {
	return (*time.Time)(t).UnmarshalText(text)
}

func (t TimeStorable) Time() time.Time {
	return time.Time(t)
}

func (ts *TimeStorable) TimePtr() *time.Time {
	if ts == nil {
		return nil
	}
	t := time.Time(*ts)
	return &t
}

var _ encoding.TextMarshaler = TimeStorable{}
var _ encoding.TextUnmarshaler = (*TimeStorable)(nil)

type HexStringStroable []byte

// MarshalText implements encoding.TextMarshaler.
func (s HexStringStroable) MarshalText() (text []byte, _ error) {
	hex.Encode(text, s)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *HexStringStroable) UnmarshalText(text []byte) (err error) {
	b := make([]byte, hex.DecodedLen(len(text)))
	_, err = hex.Decode(b, text)
	*s = b
	return
}

func (s *HexStringStroable) HexString() string {
	if s == nil {
		return ""
	}
	return hex.EncodeToString(*s)
}

var _ encoding.TextMarshaler = HexStringStroable{}
var _ encoding.TextUnmarshaler = (*HexStringStroable)(nil)

type Base64UrlStorable []byte

// MarshalText implements encoding.TextMarshaler.
func (s Base64UrlStorable) MarshalText() (text []byte, _ error) {
	text = make([]byte, base64.RawURLEncoding.EncodedLen(len(s)))
	base64.RawURLEncoding.Encode(text, s)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *Base64UrlStorable) UnmarshalText(text []byte) (err error) {
	b := make([]byte, base64.RawURLEncoding.DecodedLen(len(text)))
	_, err = base64.RawURLEncoding.Decode(b, text)
	*s = b
	return
}

func (s *Base64UrlStorable) StringPtr() *string {
	if s == nil {
		return nil
	}
	str := base64.RawStdEncoding.EncodeToString(*s)
	return &str
}

var _ encoding.TextMarshaler = Base64UrlStorable{}
var _ encoding.TextUnmarshaler = (*Base64UrlStorable)(nil)
