package cloudkey

import "crypto"

// RFC7518 3.1.  "alg" (Algorithm) Header Parameter Values for JWS

type CloudKey interface {
	KeyType() JsonWebKeyType
}

type CloudSignatureKey interface {
	CloudKey
	crypto.Signer
	KeyID() string
}

type CloudWrappingKey interface {
	CloudKey
	crypto.Decrypter
	KeyID() string
}
