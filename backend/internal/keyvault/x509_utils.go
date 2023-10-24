package kv

import (
	"encoding/asn1"
	"math/big"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

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
