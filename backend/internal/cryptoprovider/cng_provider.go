//go:build windows

package cryptoprovider

import (
	"crypto/rsa"
	"math/big"

	"github.com/microsoft/go-crypto-winnative/cng"
)

type cryptoProviderImpl struct {
}

// GenerateRSAKeyPair implements certstore.CryptoStoreProvider.
func (*cryptoProviderImpl) GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error) {
	n, e, d, p, q, dp, dq, qi, err := cng.GenerateKeyRSA(keyLength)
	if err != nil {
		return nil, err
	}
	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).SetBytes(n),
			E: int(new(big.Int).SetBytes(e).Uint64()),
		},
		D: new(big.Int).SetBytes(d),
		Primes: []*big.Int{
			new(big.Int).SetBytes(p),
			new(big.Int).SetBytes(q),
		},
		Precomputed: rsa.PrecomputedValues{
			Dp:   new(big.Int).SetBytes(dp),
			Dq:   new(big.Int).SetBytes(dq),
			Qinv: new(big.Int).SetBytes(qi),
		},
	}, nil
}

var _ CryptoProvider = (*cryptoProviderImpl)(nil)

func newCryptoProvider() *cryptoProviderImpl {
	return &cryptoProviderImpl{}
}
