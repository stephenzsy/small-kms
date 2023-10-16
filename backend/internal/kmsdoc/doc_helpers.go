package kmsdoc

import (
	"encoding"
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
	text = make([]byte, hex.EncodedLen(len(s)))
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
