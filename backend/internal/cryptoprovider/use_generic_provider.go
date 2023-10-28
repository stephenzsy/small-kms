//go:build linux

package cryptoprovider

func newCryptoProvider() CryptoProvider {
	return NewGenericCryptoProvider()
}
