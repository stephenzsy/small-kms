package cloudkey

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"math/big"
)

// RFC7518 6.1.1.  "alg" (Algorithm) Parameter Values for JWS
type JsonWebKeyType string

const (
	KeyTypeRSA JsonWebKeyType = "RSA"
	KeyTypeEC  JsonWebKeyType = "EC"
	KeyTypeOct JsonWebKeyType = "oct"
)

type JsonWebKeyCurveName string

const (
	CurveNameP256 JsonWebKeyCurveName = "P-256"
	CurveNameP384 JsonWebKeyCurveName = "P-384"
	CurveNameP521 JsonWebKeyCurveName = "P-521"
)

// RFC7517 4.3. "key_ops" (Key Operations) Parameter Values for JWK
type JsonWebKeyOperation string

const (
	JsonWebKeyOperationSign       = "sign"
	JsonWebKeyOperationVerify     = "verify"
	JsonWebKeyOperationEncrypt    = "encrypt"
	JsonWebKeyOperationDecrypt    = "decrypt"
	JsonWebKeyOperationWrapKey    = "wrapKey"
	JsonWebKeyOperationUnwrapKey  = "unwrapKey"
	JsonWebKeyOperationDeriveKey  = "deriveKey"
	JsonWebKeyOperationDeriveBits = "deriveBits"
)

type JsonWebKeyBase struct {
	KeyType          JsonWebKeyType               `json:"kty"`                // RFC7517 4.1. "kty" (Key Type) Parameter Values for JWK
	KeyID            string                       `json:"kid,omitempty"`      // RFC7517 4.5. "kid" (Key ID) Parameter
	Curve            JsonWebKeyCurveName          `json:"crv,omitempty"`      // RFC7518 6.2.1.1. "crv" (Curve) Parameter
	N                Base64RawURLEncodableBytes   `json:"n,omitempty"`        // RFC7518 6.3.1.1. "n" (Modulus) Parameter
	E                Base64RawURLEncodableBytes   `json:"e,omitempty"`        // RFC7518 6.3.1.2. "e" (Exponent) Parameter
	D                Base64RawURLEncodableBytes   `json:"d,omitempty"`        // RFC7518 6.3.2.1. "d" (Private Exponent) Parameter, or RFC7518 6.2.2.1. "d" (ECC Private Key) Parameter
	P                Base64RawURLEncodableBytes   `json:"p,omitempty"`        // RFC7518 6.3.2.2. "p" (First Prime Factor) Parameter
	Q                Base64RawURLEncodableBytes   `json:"q,omitempty"`        // RFC7518 6.3.3.3. "q" (Second Prime Factor) Parameter
	DP               Base64RawURLEncodableBytes   `json:"dp,omitempty"`       // RFC7518 6.3.3.4. "dp" (First Factor CRT Exponent) Parameter
	DQ               Base64RawURLEncodableBytes   `json:"dq,omitempty"`       // RFC7518 6.3.3.5. "dq" (Second Factor CRT Exponent) Parameter
	QI               Base64RawURLEncodableBytes   `json:"qi,omitempty"`       // RFC7518 6.3.3.6. "qi" (First CRT Coefficient) Parameter
	X                Base64RawURLEncodableBytes   `json:"x,omitempty"`        // RFC7518 6.2.1.2. "x" (X Coordinate) Parameter
	Y                Base64RawURLEncodableBytes   `json:"y,omitempty"`        // RFC7518 6.2.1.3. "y" (Y Coordinate) Parameter
	KeyOperations    []JsonWebKeyOperation        `json:"key_ops,omitempty"`  // RFC7517 4.3. "key_ops" (Key Operations) Parameter Values for JWK
	ThumbprintSHA1   Base64RawURLEncodableBytes   `json:"x5t,omitempty"`      // RFC7517 4.8. "x5t" (X.509 Certificate SHA-1 Thumbprint) Parameter
	ThumbprintSHA256 Base64RawURLEncodableBytes   `json:"x5t#S256,omitempty"` // RFC7517 4.9. "x5t#S256" (X.509 Certificate SHA-256 Thumbprint) Parameter
	CertificateChain []Base64RawURLEncodableBytes `json:"x5c,omitempty"`      // RFC7517 4.7. "x5c" (X.509 Certificate Chain) Parameter

	cachedPublicKey  crypto.PublicKey
	cachedPrivateKey crypto.PrivateKey
}

type JsonWebKey[TAlg JsonWebSignatureAlgorithm | JsonWebKeyEncryptionAlgorithm] struct {
	JsonWebKeyBase
	Alg TAlg `json:"alg,omitempty"` // RFC7517 4.4. "alg" (Algorithm) Header Parameter Values for JWS
}

type JsonWebSignatureKey = JsonWebKey[JsonWebSignatureAlgorithm]

func (jwk *JsonWebKey[T]) PublicKey() crypto.PublicKey {
	if jwk.cachedPublicKey != nil {
		return jwk.cachedPublicKey
	}

	switch jwk.KeyType {
	case KeyTypeRSA:
		return jwk.rsaPublicKey()
	case KeyTypeEC:
		return jwk.ecdsaPublicKey()
	}
	return jwk.cachedPublicKey
}

func (jwk *JsonWebKey[T]) rsaPublicKey() *rsa.PublicKey {
	if (jwk.cachedPublicKey) != nil {
		return jwk.cachedPublicKey.(*rsa.PublicKey)
	}
	jwk.cachedPublicKey = &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(jwk.N),
		E: int(big.NewInt(0).SetBytes(jwk.E).Int64()),
	}
	return jwk.cachedPublicKey.(*rsa.PublicKey)
}

func (jwk *JsonWebKey[T]) ecdsaPublicKey() *ecdsa.PublicKey {
	if (jwk.cachedPublicKey) != nil {
		return jwk.cachedPublicKey.(*ecdsa.PublicKey)
	}
	var crv elliptic.Curve
	switch jwk.Curve {
	case CurveNameP256:
		crv = elliptic.P256()
	case CurveNameP384:
		crv = elliptic.P384()
	case CurveNameP521:
		crv = elliptic.P521()
	default:
		return nil
	}
	jwk.cachedPublicKey = &ecdsa.PublicKey{
		Curve: crv,
		X:     big.NewInt(0).SetBytes(jwk.X),
		Y:     big.NewInt(0).SetBytes(jwk.Y),
	}
	return jwk.cachedPublicKey.(*ecdsa.PublicKey)
}

// Cloud keys typically don't have retrieveable private key
func (jwk *JsonWebKey[T]) PrivateKey() crypto.PrivateKey {
	if jwk.cachedPrivateKey != nil {
		return jwk.cachedPrivateKey
	}

	switch jwk.KeyType {
	case KeyTypeRSA:
		jwk.cachedPrivateKey = &rsa.PrivateKey{
			PublicKey: *jwk.rsaPublicKey(),
			D:         big.NewInt(0).SetBytes(jwk.D),
			Primes:    []*big.Int{big.NewInt(0).SetBytes(jwk.P), big.NewInt(0).SetBytes(jwk.Q)},
			Precomputed: rsa.PrecomputedValues{
				Dp:   big.NewInt(0).SetBytes(jwk.DP),
				Dq:   big.NewInt(0).SetBytes(jwk.DQ),
				Qinv: big.NewInt(0).SetBytes(jwk.QI),
			},
		}
		return jwk.cachedPrivateKey
	case KeyTypeEC:
		jwk.cachedPrivateKey = &ecdsa.PrivateKey{
			PublicKey: *jwk.ecdsaPublicKey(),
			D:         big.NewInt(0).SetBytes(jwk.D),
		}
	}

	return jwk.cachedPrivateKey
}
