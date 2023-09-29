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
	"math/big"
	"net"
	"net/url"
	"slices"
	"strings"
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
		slices.Sort(trimmed)
		*ptr = trimmed
	}
}

func sanitizeStringArrayWithParse[D any](ptr *[]string, parsed *[]D, parse func(string) (D, error)) error {
	if *ptr == nil {
		return nil
	}
	if len(*ptr) <= 0 {
		ptr = nil
		return nil
	}
	v := make([]string, 0, len(*ptr))
	vp := make([]D, 0, len(*ptr))
	for _, s := range *ptr {
		s = strings.TrimSpace(s)
		if vs, err := parse(s); err == nil {
			vp = append(vp, vs)
			v = append(v, s)
		} else {
			return err
		}
	}
	if len(v) > 0 {
		slices.Sort(v)
		*ptr = v
		*parsed = vp
	} else {
		*ptr = nil
		*parsed = nil
	}
	return nil
}

type SANsSanitized struct {
	CertificateSubjectAlternativeNames
	parsedIPAddresses []net.IP
	parsedURLs        []*url.URL
}

func sanitizeSANs(san *CertificateSubjectAlternativeNames) (*SANsSanitized, error) {
	if san == nil {
		return nil, nil
	}
	var parsedIPAddresses []net.IP
	var parsedURLs []*url.URL
	sanitizeStringArray(&san.DNSNames)
	sanitizeStringArray(&san.EmailAddresses)
	err := sanitizeStringArrayWithParse[net.IP](&san.IPAddresses, &parsedIPAddresses, func(s string) (ip net.IP, err error) {
		ip = net.ParseIP(s)
		if ip == nil {
			err = fmt.Errorf("invalid ip address: %s", s)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	err = sanitizeStringArrayWithParse[*url.URL](&san.URIs, &parsedURLs, url.Parse)
	if err != nil {
		return nil, err
	}
	if len(san.DNSNames)+
		len(san.EmailAddresses)+
		len(san.IPAddresses)+
		len(san.URIs) == 0 {
		return nil, nil
	}
	return &SANsSanitized{
		CertificateSubjectAlternativeNames: *san,
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

func (sans *SANsSanitized) populateCertificate(cert *x509.Certificate) {
	if sans == nil {
		return
	}
	cert.DNSNames = sans.DNSNames
	cert.EmailAddresses = sans.EmailAddresses
	cert.IPAddresses = sans.parsedIPAddresses
	cert.URIs = sans.parsedURLs
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
			//	p.E = ToPtr(base64.URLEncoding.EncodeToString(big.NewInt(int64(rsaPublicKey.E)).Bytes()))
			//	p.N = ToPtr(base64.URLEncoding.EncodeToString(rsaPublicKey.N.Bytes()))
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
			//	p.X = ToPtr(base64.URLEncoding.EncodeToString(ecdsaPublicKey.X.Bytes()))
			//	p.Y = ToPtr(base64.URLEncoding.EncodeToString(ecdsaPublicKey.Y.Bytes()))
		}
	}
	thumbprinterSha1 := sha1.New().Sum(c.Raw)
	thumbprinterSha256 := sha256.New().Sum(c.Raw)

	p.CertificateThumbprint = ToPtr(base64.URLEncoding.EncodeToString(thumbprinterSha1))
	p.CertificateThumbprintSHA256 = ToPtr(base64.URLEncoding.EncodeToString(thumbprinterSha256))

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
