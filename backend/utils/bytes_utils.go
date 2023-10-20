package utils

import (
	"crypto/sha1"
	"encoding"
	"encoding/hex"
)

type CertificateFingerprintSHA1 [sha1.Size]byte

func (s CertificateFingerprintSHA1) MarshalText() (text []byte, _ error) {
	text = make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(text, s[:])
	return
}

func (s *CertificateFingerprintSHA1) UnmarshalText(text []byte) (err error) {
	bs := make([]byte, hex.DecodedLen(len(text)))
	if l, err := hex.Decode(bs, text); err != nil {
		return err
	} else if l == sha1.Size {
		*s = ([sha1.Size]byte)(bs)
	}
	return
}

func (s *CertificateFingerprintSHA1) HexString() string {
	if s == nil {
		return ""
	}
	return hex.EncodeToString(s[:])
}

var _ encoding.TextMarshaler = CertificateFingerprintSHA1{}
var _ encoding.TextUnmarshaler = (*CertificateFingerprintSHA1)(nil)
