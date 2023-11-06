package cert

import (
	"crypto/x509/pkix"
	"encoding"
	"encoding/hex"
	"io"

	"github.com/stephenzsy/small-kms/backend/base"
)

type ID = base.ID

type (
	certPolicyRefComposed struct {
		base.ResourceReference
		CertPolicyRefFields
	}

	certPolicyComposed struct {
		CertPolicyRef
		CertPolicyFields
	}

	certificateRefComposed struct {
		base.ResourceReference
		CertificateRefFields
	}

	certificateComposed struct {
		CertificateRef
		CertificateFields
	}
)

func (s *CertificateSubject) ToPkixName() pkix.Name {
	return pkix.Name{
		CommonName: s.CommonName,
	}
}

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

var _ encoding.TextMarshaler = HexDigest{}
var _ encoding.TextUnmarshaler = (*HexDigest)(nil)

func (sans *SubjectAlternativeNames) WriteToDigest(w io.Writer) (s int, err error) {
	if sans == nil {
		return 0, nil
	}
	for _, san := range sans.DNSNames {
		if c, err := w.Write([]byte(san)); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	for _, san := range sans.Emails {
		if c, err := w.Write([]byte(san)); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	for _, san := range sans.IPAddresses {
		if c, err := w.Write(san); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	return s, nil
}
