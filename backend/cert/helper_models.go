package cert

import (
	"encoding"
	"math/big"

	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
)

type SerialNumberStorable kmsdoc.HexStringStroable

// MarshalText implements encoding.TextMarshaler.
func (s SerialNumberStorable) MarshalText() (text []byte, _ error) {
	return kmsdoc.HexStringStroable(s).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *SerialNumberStorable) UnmarshalText(text []byte) (err error) {
	return (*kmsdoc.HexStringStroable)(s).UnmarshalText(text)
}

func (s SerialNumberStorable) BigInt() *big.Int {
	return new(big.Int).SetBytes(s)
}

var _ encoding.TextMarshaler = SerialNumberStorable{}
var _ encoding.TextUnmarshaler = (*SerialNumberStorable)(nil)
