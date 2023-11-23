package cloudkeyaz

import (
	"encoding/asn1"
	"math/big"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	i "github.com/stephenzsy/small-kms/backend/cloud/key"
)

func ExtractKeyVaultName(keyvaultEndpoing string) string {
	if parsed, err := url.Parse(keyvaultEndpoing); err == nil {
		return strings.Split(parsed.Host, ".")[0]
	}
	return ""
}

func newSigningJWKFromKeyVaultKey(kvKey *azkeys.JSONWebKey) *i.JsonWebKey {
	keyID := string(*kvKey.KID)
	r := &i.JsonWebKey{
		KeyType: ktyReverseMapping[*kvKey.Kty],
		KeyID:   keyID,
		X:       kvKey.X,
		Y:       kvKey.Y,
		N:       kvKey.N,
		E:       kvKey.E,
	}
	if kvKey.Crv != nil {
		crv := crvReverseMapping[*kvKey.Crv]
		r.Curve = crv
	}
	return r
}

type esSignature struct {
	R *big.Int
	S *big.Int
}

func toX509Signature(signedDigest []byte, sigAlg azkeys.SignatureAlgorithm) ([]byte, error) {
	var n uint32
	switch sigAlg {
	case azkeys.SignatureAlgorithmES256:
		n = 32
	case azkeys.SignatureAlgorithmES384:
		n = 48
	case azkeys.SignatureAlgorithmES512:
		n = 66
	default:
		return signedDigest, nil
	}
	sig := esSignature{
		R: new(big.Int).SetBytes(signedDigest[:n]),
		S: new(big.Int).SetBytes(signedDigest[n:]),
	}
	return asn1.Marshal(sig)
}
