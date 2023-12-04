package resdoc

import (
	"encoding"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
)

type NumericDate = jwt.NumericDate

type HexBytes []byte

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *HexBytes) UnmarshalText(text []byte) error {
	sl := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(sl, text)
	if err != nil {
		return err
	}
	*s = sl
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (s HexBytes) MarshalText() (text []byte, err error) {
	text = make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(text, s)
	return
}

func (s HexBytes) String() string {
	return hex.EncodeToString(s)
}

var _ encoding.TextMarshaler = HexBytes{}
var _ encoding.TextUnmarshaler = (*HexBytes)(nil)
