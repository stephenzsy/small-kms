//go:build linux && cgo

package cryptoprovider

import (
	"crypto/rsa"
	"math/big"

	"github.com/microsoft/go-crypto-openssl/openssl"
	"github.com/microsoft/go-crypto-openssl/openssl/bbig"
)

type cryptoProviderImpl struct {
}

// GenerateRSAKeyPair implements certstore.CryptoStoreProvider.
func (*cryptoProviderImpl) GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error) {
	//vs
	n, e, d, p, q, dp, dq, qi, err := openssl.GenerateKeyRSA(keyLength)
	if err != nil {
		return nil, err
	}
	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: bbig.Dec(n),
			E: int(bbig.Dec(e).Uint64()),
		},
		D: bbig.Dec(d),
		Primes: []*big.Int{
			bbig.Dec(p),
			bbig.Dec(q),
		},
		Precomputed: rsa.PrecomputedValues{
			Dp:   bbig.Dec(dp),
			Dq:   bbig.Dec(dq),
			Qinv: bbig.Dec(qi),
		},
	}, nil
}

var _ CryptoProvider = (*cryptoProviderImpl)(nil)

func newCryptoProvider() *cryptoProviderImpl {
	return &cryptoProviderImpl{}
}
