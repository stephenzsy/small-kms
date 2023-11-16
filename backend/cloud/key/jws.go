package cloudkey

import "crypto"

type JsonWebSignatureAlgorithm string

const (
	SignatureAlgoritmNone JsonWebSignatureAlgorithm = ""

	SignatureAlgorithmHS256 JsonWebSignatureAlgorithm = "HS256"
	SignatureAlgorithmHS384 JsonWebSignatureAlgorithm = "HS384"
	SignatureAlgorithmHS512 JsonWebSignatureAlgorithm = "HS512"

	SignatureAlgorithmRS256 JsonWebSignatureAlgorithm = "RS256"
	SignatureAlgorithmRS384 JsonWebSignatureAlgorithm = "RS384"
	SignatureAlgorithmRS512 JsonWebSignatureAlgorithm = "RS512"

	SignatureAlgorithmES256 JsonWebSignatureAlgorithm = "ES256"
	SignatureAlgorithmES384 JsonWebSignatureAlgorithm = "ES384"
	SignatureAlgorithmES512 JsonWebSignatureAlgorithm = "ES512"

	SignatureAlgorithmPS256 JsonWebSignatureAlgorithm = "PS256"
	SignatureAlgorithmPS384 JsonWebSignatureAlgorithm = "PS384"
	SignatureAlgorithmPS512 JsonWebSignatureAlgorithm = "PS512"
)

var supportedAlgs = map[JsonWebSignatureAlgorithm]bool{
	SignatureAlgorithmHS256: true,
	SignatureAlgorithmHS384: true,
	SignatureAlgorithmHS512: true,
	SignatureAlgorithmRS256: true,
	SignatureAlgorithmRS384: true,
	SignatureAlgorithmRS512: true,
	SignatureAlgorithmES256: true,
	SignatureAlgorithmES384: true,
	SignatureAlgorithmES512: true,
	SignatureAlgorithmPS256: true,
	SignatureAlgorithmPS384: true,
	SignatureAlgorithmPS512: true,
}

// HashFunc implements crypto.SignerOpts.
func (alg JsonWebSignatureAlgorithm) HashFunc() crypto.Hash {
	switch alg {
	case SignatureAlgorithmHS256,
		SignatureAlgorithmRS256,
		SignatureAlgorithmES256,
		SignatureAlgorithmPS256:
		return crypto.SHA256
	case SignatureAlgorithmHS384,
		SignatureAlgorithmRS384,
		SignatureAlgorithmES384,
		SignatureAlgorithmPS384:
		return crypto.SHA384
	case SignatureAlgorithmHS512,
		SignatureAlgorithmRS512,
		SignatureAlgorithmES512,
		SignatureAlgorithmPS512:
		return crypto.SHA512
	default:
		return 0
	}
}

// HashFunc implements crypto.SignerOpts.
func (alg JsonWebSignatureAlgorithm) IsSupported() bool {
	return supportedAlgs[alg]
}

var _ crypto.SignerOpts = JsonWebSignatureAlgorithm("")
