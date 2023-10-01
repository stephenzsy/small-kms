package admin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash"
	"math/big"
	"net"
	"net/url"
	"slices"
	"strings"

	"github.com/stephenzsy/small-kms/backend/common"
)

func sanitizeStringArray(ptr *[]string) {
	if *ptr == nil {
		return
	}
	if len(*ptr) <= 0 {
		ptr = nil
		return
	}
	trimmed := make([]string, 0, len(*ptr))
	for _, s := range *ptr {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			trimmed = append(trimmed, s)
		}
	}
	if len(trimmed) > 0 {
		slices.SortFunc(trimmed, func(a, b string) int {
			return strings.Compare(strings.ToUpper(a), strings.ToUpper(b))
		})
		*ptr = trimmed
	}
}

type parsedSan[D net.IP | *url.URL] struct {
	StringValue string
	Parsed      D
	Variable    *common.CertificateFieldVar
}

func sanitizeStringArrayWithParse[D net.IP | *url.URL](ptr *[]string, parseFunc func(string) (parsedSan[D], error), sortFunc func(a, b parsedSan[D]) int) ([]parsedSan[D], error) {
	if ptr == nil || len(*ptr) == 0 {
		*ptr = nil
		return nil, nil
	}
	rslice := make([]parsedSan[D], 0, len(*ptr))
	dedup := make(map[string]bool)
	for _, s := range *ptr {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		if item, err := parseFunc(s); err != nil {
			return nil, err
		} else if !dedup[item.StringValue] {
			dedup[item.StringValue] = true
			rslice = append(rslice, item)
		}
	}
	if len(rslice) == 0 {
		*ptr = nil
		return nil, nil
	}
	slices.SortFunc(rslice, sortFunc)
	sortedStringSlice := make([]string, len(rslice))
	for i, item := range rslice {
		sortedStringSlice[i] = item.StringValue
	}
	*ptr = sortedStringSlice
	return rslice, nil
}

type SANsSanitized struct {
	CertificateSubjectAlternativeNames
	parsedIPAddresses []parsedSan[net.IP]
	parsedURLs        []parsedSan[*url.URL]
}

func parseSanIPAddr(s string) (p parsedSan[net.IP], err error) {
	p.Parsed = net.ParseIP(s)
	if p.Parsed == nil {
		err = fmt.Errorf("invalid ip address: %s", s)
		return
	}
	p.StringValue = p.Parsed.String()
	return
}

func parseSanURL(s string) (p parsedSan[*url.URL], err error) {
	p.Parsed, err = url.Parse(s)
	if err != nil {
		return
	}
	p.StringValue = p.Parsed.String()
	return
}

func parseSanURLWithVariable(s string) (p parsedSan[*url.URL], err error) {
	v, isVar, err := validateCertFieldForVariable(&s)
	if err != nil {
		return
	}
	if isVar {
		p.Parsed = nil
		p.StringValue = s
		p.Variable = v
		return
	}
	return parseSanURL(s)
}

func sanitizeSANs(sans *CertificateSubjectAlternativeNames) (*SANsSanitized, error) {
	if sans == nil {
		return nil, nil
	}

	sanitizeStringArray(&sans.DNSNames)
	sanitizeStringArray(&sans.EmailAddresses)

	parsedIPAddresses, err := sanitizeStringArrayWithParse[net.IP](&sans.IPAddresses,
		parseSanIPAddr, func(a, b parsedSan[net.IP]) int {
			// ipv4 before ipv6
			lenDiff := len(a.Parsed) - len(b.Parsed)
			if lenDiff != 0 {
				return lenDiff
			}
			return slices.Compare(a.Parsed, b.Parsed)
		})

	if err != nil {
		return nil, err
	}
	parsedURLs, err := sanitizeStringArrayWithParse[*url.URL](&sans.URIs,
		parseSanURLWithVariable, func(a, b parsedSan[*url.URL]) int {
			return strings.Compare(a.StringValue, b.StringValue)
		})
	if err != nil {
		return nil, err
	}
	if len(sans.DNSNames)+
		len(sans.EmailAddresses)+
		len(sans.IPAddresses)+
		len(sans.URIs) == 0 {
		return nil, nil
	}
	return &SANsSanitized{
		CertificateSubjectAlternativeNames: *sans,
		parsedIPAddresses:                  parsedIPAddresses,
		parsedURLs:                         parsedURLs,
	}, nil
}

func (a *SANsSanitized) Equals(b *SANsSanitized) bool {
	if a == nil {
		return b == nil
	}
	return slices.Equal(a.DNSNames, b.DNSNames) &&
		slices.Equal(a.EmailAddresses, b.EmailAddresses) &&
		slices.Equal(a.IPAddresses, b.IPAddresses) &&
		slices.Equal(a.URIs, b.URIs)
}

func (sans *SANsSanitized) populateCertificate(cert *x509.Certificate, variableValues map[string]string) (err error) {
	if sans == nil {
		return
	}
	cert.DNSNames = sans.DNSNames
	cert.EmailAddresses = sans.EmailAddresses
	// reparse
	cert.IPAddresses, err = common.SliceMapWithError(sans.IPAddresses, func(s string) (net.IP, error) {
		p, err := parseSanIPAddr(s)
		return p.Parsed, err
	})
	if err != nil {
		return
	}
	cert.URIs, err = common.SliceMapWithError(sans.URIs, func(s string) (*url.URL, error) {
		p, err := parseSanURL(s)
		if err != nil {
			return nil, err
		}
		if p.Variable != nil {
			// make subsitution, and redo parsing
			s, err = p.Variable.Substitute(variableValues)
			if err != nil {
				return nil, err
			}
			p, err = parseSanURL(s)
		}
		return p.Parsed, err
	})
	return
}

func (p *JwkProperties) populateBriefFromCertificate(c *x509.Certificate) {
	if p == nil || c == nil {
		return
	}
	switch c.SignatureAlgorithm {
	case x509.SHA256WithRSA:
		p.Alg = ToPtr(AlgRS256)
	case x509.SHA384WithRSA:
		p.Alg = ToPtr(AlgRS384)
	case x509.SHA512WithRSA:
		p.Alg = ToPtr(AlgRS512)
	case x509.ECDSAWithSHA256:
		p.Alg = ToPtr(AlgES256)
	case x509.ECDSAWithSHA384:
		p.Alg = ToPtr(AlgES384)
	}

	switch c.PublicKeyAlgorithm {
	case x509.RSA:
		p.Kty = KeyTypeRSA
		if rsaPublicKey, ok := c.PublicKey.(*rsa.PublicKey); ok {
			p.KeySize = ToPtr(KeySize(rsaPublicKey.N.BitLen()))
		}
	case x509.ECDSA:
		p.Kty = KeyTypeEC
		if ecdsaPublicKey, ok := c.PublicKey.(*ecdsa.PublicKey); ok {
			curveName := ecdsaPublicKey.Curve.Params().Name
			switch curveName {
			case elliptic.P256().Params().Name:
				p.Crv = ToPtr(CurveNameP256)
			case elliptic.P384().Params().Name:
				p.Crv = ToPtr(CurveNameP384)
			}
		}
	}
	thumbprinterSha1 := getThumbprint(sha1.New(), c.Raw)
	thumbprinterSha256 := getThumbprint(sha256.New(), c.Raw)

	p.CertificateThumbprint = ToPtr(base64.URLEncoding.EncodeToString(thumbprinterSha1))
	p.CertificateThumbprintSHA256 = ToPtr(base64.URLEncoding.EncodeToString(thumbprinterSha256))
}

func getThumbprint(h hash.Hash, b []byte) []byte {
	h.Write(b)
	return h.Sum(nil)
}

func (p *JwkProperties) populateCertsFromPemBlob(pemBlob []byte) error {
	if p == nil || len(pemBlob) <= 0 {
		return nil
	}

	pemBlock, restPem := pem.Decode(pemBlob)
	c, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}
	switch p.Kty {
	case KeyTypeRSA:
		if rsaPublicKey, ok := c.PublicKey.(*rsa.PublicKey); ok {
			p.E = ToPtr(base64.URLEncoding.EncodeToString(big.NewInt(int64(rsaPublicKey.E)).Bytes()))
			p.N = ToPtr(base64.URLEncoding.EncodeToString(rsaPublicKey.N.Bytes()))
		}
	case KeyTypeEC:
		if ecdsaPublicKey, ok := c.PublicKey.(*ecdsa.PublicKey); ok {
			p.X = ToPtr(base64.URLEncoding.EncodeToString(ecdsaPublicKey.X.Bytes()))
			p.Y = ToPtr(base64.URLEncoding.EncodeToString(ecdsaPublicKey.Y.Bytes()))
		}
	}
	p.CertificateChain = make([]string, 1, 3)
	p.CertificateChain[0] = base64.URLEncoding.EncodeToString(c.Raw)
	for block, rest := pem.Decode(restPem); block != nil; block, rest = pem.Decode(rest) {
		p.CertificateChain = append(p.CertificateChain, base64.URLEncoding.EncodeToString(block.Bytes))
	}
	return nil
}
