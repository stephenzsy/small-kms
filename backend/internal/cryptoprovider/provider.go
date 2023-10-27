// go:build windows || (linux && cgo)
package cryptoprovider

import (
	"crypto/rsa"
)

type CryptoProvider interface {
	GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error)
}

func NewCryptoProvider() (CryptoProvider, error) {
	return newCryptoProvider(), nil
}
