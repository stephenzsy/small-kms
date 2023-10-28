package cryptoprovider

import (
	"crypto/rand"
	"crypto/rsa"
)

type genericCryptoProviderImpl struct {
}

// GenerateRSAKeyPair implements certstore.CryptoStoreProvider.
func (*genericCryptoProviderImpl) GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keyLength)

}

var _ CryptoProvider = (*genericCryptoProviderImpl)(nil)

func NewGenericCryptoProvider() CryptoProvider {
	return &genericCryptoProviderImpl{}
}
