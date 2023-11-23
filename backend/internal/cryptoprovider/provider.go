// go:build windows || (linux && cgo)
package cryptoprovider

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
)

type CryptoProvider interface {
	GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error)
	GenerateECDSAKeyPair(c elliptic.Curve) (*ecdsa.PrivateKey, error)
}

func NewCryptoProvider() (CryptoProvider, error) {
	return newCryptoProvider(), nil
}
