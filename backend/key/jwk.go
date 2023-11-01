package key

import (
	"crypto/rsa"
	"fmt"
	"math/big"

	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
)

func (jwk *JsonWebKey) AsRsaPubicKey() (*rsa.PublicKey, error) {
	if jwk.Kty != cloudkey.KeyTypeRSA || jwk.E == nil || jwk.N == nil {
		return nil, fmt.Errorf("invalid public key type")
	}
	return &rsa.PublicKey{
		E: int(big.NewInt(int64(0)).SetBytes(jwk.E).Int64()),
		N: big.NewInt(int64(0)).SetBytes(jwk.N),
	}, nil
}
