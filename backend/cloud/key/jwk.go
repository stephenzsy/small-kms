package cloudkey

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

type Base64RawURLEncodableBytes []byte

// MarshalText implements encoding.TextMarshaler.
func (b Base64RawURLEncodableBytes) MarshalText() (text []byte, err error) {
	text = make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(text, b)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Base64RawURLEncodableBytes) UnmarshalText(text []byte) error {
	*b = make([]byte, base64.RawURLEncoding.DecodedLen(len(text)))
	_, err := base64.RawURLEncoding.Decode(*b, text)
	return err
}

func (b Base64RawURLEncodableBytes) HexString() string {
	return hex.EncodeToString(b)
}

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
	JosnWebKeyOperationUnwrapKey  = "unwrapKey"
	JsonWebKeyOperationDeriveKey  = "deriveKey"
	JsonWebKeyOperationDeriveBits = "deriveBits"
)

type JsonWebKey[TAlg JsonWebSignatureAlgorithm] struct {
	KeyType       JsonWebKeyType             `json:"kty"`               // RFC7517 4.1. "kty" (Key Type) Parameter Values for JWK
	Alg           TAlg                       `json:"alg,omitempty"`     // RFC7517 4.4. "alg" (Algorithm) Header Parameter Values for JWS
	KeyID         string                     `json:"kid,omitempty"`     // RFC7517 4.5. "kid" (Key ID) Parameter
	CurveName     JsonWebKeyCurveName        `json:"crv,omitempty"`     // RFC7518 6.2.1.1. "crv" (Curve) Parameter
	N             Base64RawURLEncodableBytes `json:"n,omitempty"`       // RFC7518 6.3.1.1. "n" (Modulus) Parameter
	E             Base64RawURLEncodableBytes `json:"e,omitempty"`       // RFC7518 6.3.1.2. "e" (Exponent) Parameter
	X             Base64RawURLEncodableBytes `json:"x,omitempty"`       // RFC7518 6.2.1.2. "x" (X Coordinate) Parameter
	Y             Base64RawURLEncodableBytes `json:"y,omitempty"`       // RFC7518 6.2.1.3. "y" (Y Coordinate) Parameter
	KeyOperations []string                   `json:"key_ops,omitempty"` // RFC7517 4.3. "key_ops" (Key Operations) Parameter Values for JWK

	cachedPublicKey crypto.PublicKey
}

type JsonWebSignatureKey = JsonWebKey[JsonWebSignatureAlgorithm]

func (jwk *JsonWebKey[T]) PublicKey() crypto.PublicKey {
	if jwk.cachedPublicKey != nil {
		return jwk.cachedPublicKey
	}

	switch jwk.KeyType {
	case KeyTypeRSA:
		jwk.cachedPublicKey = &rsa.PublicKey{
			N: big.NewInt(0).SetBytes(jwk.N),
			E: int(big.NewInt(0).SetBytes(jwk.E).Int64()),
		}
		return jwk.cachedPublicKey
	case KeyTypeEC:
		var crv elliptic.Curve
		switch jwk.CurveName {
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
	}
	return jwk.cachedPublicKey
}
