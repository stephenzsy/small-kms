package cryptoprovider

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
)

type genericCryptoProviderImpl struct {
}

// GenerateECDSAKeyPair implements CryptoProvider.
func (*genericCryptoProviderImpl) GenerateECDSAKeyPair(c elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(c, rand.Reader)
}

// GenerateRSAKeyPair implements certstore.CryptoStoreProvider.
func (*genericCryptoProviderImpl) GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keyLength)

}

var _ CryptoProvider = (*genericCryptoProviderImpl)(nil)

func NewGenericCryptoProvider() CryptoProvider {
	return &genericCryptoProviderImpl{}
}
