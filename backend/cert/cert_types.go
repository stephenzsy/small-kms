package cert

import (
	"context"
	"crypto/x509/pkix"
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

type HexDigest = base.HexDigest

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

type internalContextKey int

const (
	groupMemberOfContextKey internalContextKey = iota
	groupMemberGraphObjectContextKey
	selfGraphObjectContextKey
)

func (s *CertificateSubject) processTemplate(c context.Context) (processed CertificateSubject, err error) {
	if cn, err := ProcessTemplate(c, "subjectCN", s.CommonName); err != nil {
		return processed, err
	} else {
		processed.CommonName = cn
	}
	return
}
