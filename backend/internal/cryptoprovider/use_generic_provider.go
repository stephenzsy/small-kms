//go:build linux || darwin

package cryptoprovider

func newCryptoProvider() CryptoProvider {
	return NewGenericCryptoProvider()
}
