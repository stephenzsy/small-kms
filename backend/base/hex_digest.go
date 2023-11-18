package base

import (
	"encoding"
	"encoding/hex"
)

type HexDigest []byte

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *HexDigest) UnmarshalText(text []byte) error {
	sl := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(sl, text)
	if err != nil {
		return err
	}
	*s = sl
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (s HexDigest) MarshalText() (text []byte, err error) {
	text = make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(text, s)
	return
}

func (s HexDigest) String() string {
	return hex.EncodeToString(s)
}

var _ encoding.TextMarshaler = HexDigest{}
var _ encoding.TextUnmarshaler = (*HexDigest)(nil)
