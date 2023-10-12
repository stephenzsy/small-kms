package shared

import (
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

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

func (alg JwkAlg) ToX509SignatureAlgorithm() x509.SignatureAlgorithm {
	switch alg {
	case AlgRS256:
		return x509.SHA256WithRSA
	case AlgRS384:
		return x509.SHA384WithRSA
	case AlgRS512:
		return x509.SHA512WithRSA
	case AlgES256:
		return x509.ECDSAWithSHA256
	case AlgES384:
		return x509.ECDSAWithSHA384
	}
	return x509.UnknownSignatureAlgorithm
}
