package pkcs12utils

import (
	"crypto/x509"

	"software.sslmate.com/src/go-pkcs12"
)

func ConvertPKCS12(privateKey any, certChain []*x509.Certificate, password string, legacy bool) ([]byte, error) {
	encoder := pkcs12.Modern
	if legacy {
		encoder = pkcs12.Legacy
	}
	return encoder.Encode(privateKey, certChain[0], certChain[1:], password)
}
