package models

import "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"

func (alg JwkAlg) ToAzKeysSignatureAlgorithm() azkeys.SignatureAlgorithm {
	switch alg {
	case AlgRS256:
		return azkeys.SignatureAlgorithmRS256
	case AlgRS384:
		return azkeys.SignatureAlgorithmRS384
	case AlgRS512:
		return azkeys.SignatureAlgorithmRS512
	case AlgES256:
		return azkeys.SignatureAlgorithmES256
	case AlgES384:
		return azkeys.SignatureAlgorithmES384
	}
	return azkeys.SignatureAlgorithm("")
}
