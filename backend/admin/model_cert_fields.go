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
	"hash"
	"math/big"
	"net/url"
	"slices"
	"strings"

	"github.com/stephenzsy/small-kms/backend/utils"
)

func sanitizeStringArray(ptr *[]string) {
	if *ptr == nil {
		return
	}
	if len(*ptr) <= 0 {
		ptr = nil
		return
	}
	for i, s := range *ptr {
		(*ptr)[i] = strings.TrimSpace(s)
	}
	*ptr = utils.FilterSlice(*ptr, func(s string) bool { return len(s) > 0 })
	slices.SortFunc(*ptr, func(a, b string) int {
		return strings.Compare(strings.ToUpper(a), strings.ToUpper(b))
	})
	*ptr = slices.Compact(*ptr)
	if len(*ptr) == 0 {
		*ptr = nil
	}
}

func sanitizeSANs(sans *CertificateSubjectAlternativeNames) *CertificateSubjectAlternativeNames {
	if sans == nil {
		return nil
	}

	sanitizeStringArray(&sans.EmailAddresses)
	sanitizeStringArray(&sans.URIs)
	if sans.EmailAddresses == nil && sans.URIs == nil {
		return nil
	}
	return sans
}

func (sans *CertificateSubjectAlternativeNames) populateCertificate(cert *x509.Certificate, data *TemplateVarData) {
	if sans == nil {
		return
	}

	cert.EmailAddresses = utils.NilIfZeroLen(
		utils.FilterSlice(
			utils.MapSlices(sans.EmailAddresses, func(emailAddrStr string) string {
				return processTemplate(emailAddrStr, data)
			}),
			func(emailAddrStr string) bool {
				return len(emailAddrStr) > 0
			}))

	cert.URIs = utils.NilIfZeroLen[*url.URL](
		utils.FilterSlice[*url.URL](
			utils.MapSlices(sans.URIs, func(uriStr string) *url.URL {
				uriStr = processTemplate(uriStr, data)
				if len(uriStr) > 0 {
					uri, _ := url.Parse(uriStr)
					return uri
				}
				return nil
			}),
			func(uri *url.URL) bool {
				return uri != nil
			}))
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
