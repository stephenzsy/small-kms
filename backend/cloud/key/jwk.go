package cloudkey

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"io"
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

type JsonWebKey struct {
	KeyType          JsonWebKeyType               `json:"kty"`                // RFC7517 4.1. "kty" (Key Type) Parameter Values for JWK
	Alg              string                       `json:"alg"`                // RFC7517 4.4. "alg" (Algorithm) Header Parameter Values for JWS
	KeyID            string                       `json:"kid,omitempty"`      // RFC7517 4.5. "kid" (Key ID) Parameter
	Curve            JsonWebKeyCurveName          `json:"crv,omitempty"`      // RFC7518 6.2.1.1. "crv" (Curve) Parameter
	N                Base64RawURLEncodableBytes   `json:"n,omitempty"`        // RFC7518 6.3.1.1. "n" (Modulus) Parameter
	E                Base64RawURLEncodableBytes   `json:"e,omitempty"`        // RFC7518 6.3.1.2. "e" (Exponent) Parameter
	D                Base64RawURLEncodableBytes   `json:"d,omitempty"`        // RFC7518 6.3.2.1. "d" (Private Exponent) Parameter, or RFC7518 6.2.2.1. "d" (ECC Private Key) Parameter
	P                Base64RawURLEncodableBytes   `json:"p,omitempty"`        // RFC7518 6.3.2.2. "p" (First Prime Factor) Parameter
	Q                Base64RawURLEncodableBytes   `json:"q,omitempty"`        // RFC7518 6.3.3.3. "q" (Second Prime Factor) Parameter
	Dp               Base64RawURLEncodableBytes   `json:"dp,omitempty"`       // RFC7518 6.3.3.4. "dp" (First Factor CRT Exponent) Parameter
	Dq               Base64RawURLEncodableBytes   `json:"dq,omitempty"`       // RFC7518 6.3.3.5. "dq" (Second Factor CRT Exponent) Parameter
	Qinv             Base64RawURLEncodableBytes   `json:"qi,omitempty"`       // RFC7518 6.3.3.6. "qi" (First CRT Coefficient) Parameter
	X                Base64RawURLEncodableBytes   `json:"x,omitempty"`        // RFC7518 6.2.1.2. "x" (X Coordinate) Parameter
	Y                Base64RawURLEncodableBytes   `json:"y,omitempty"`        // RFC7518 6.2.1.3. "y" (Y Coordinate) Parameter
	KeyOperations    []JsonWebKeyOperation        `json:"key_ops,omitempty"`  // RFC7517 4.3. "key_ops" (Key Operations) Parameter Values for JWK
	ThumbprintSHA1   Base64RawURLEncodableBytes   `json:"x5t,omitempty"`      // RFC7517 4.8. "x5t" (X.509 Certificate SHA-1 Thumbprint) Parameter
	ThumbprintSHA256 Base64RawURLEncodableBytes   `json:"x5t#S256,omitempty"` // RFC7517 4.9. "x5t#S256" (X.509 Certificate SHA-256 Thumbprint) Parameter
	CertificateChain []Base64RawURLEncodableBytes `json:"x5c,omitempty"`      // RFC7517 4.7. "x5c" (X.509 Certificate Chain) Parameter

	cachedPublicKey  crypto.PublicKey
	cachedPrivateKey crypto.PrivateKey
}

func (jwk *JsonWebKey) Digest(w io.Writer) {
	w.Write([]byte(jwk.KeyType))
	w.Write([]byte(jwk.Curve))
	w.Write(jwk.N)
	w.Write(jwk.E)
	w.Write(jwk.X)
	w.Write(jwk.Y)
	w.Write([]byte(jwk.KeyID))
	w.Write(jwk.ThumbprintSHA1)
	w.Write(jwk.ThumbprintSHA256)
	for _, v := range jwk.CertificateChain {
		w.Write(v)
	}
	for _, v := range jwk.KeyOperations {
		w.Write([]byte(v))
	}
}

func (jwk *JsonWebKey) PublicJWK() *JsonWebKey {
	return &JsonWebKey{
		KeyType:          jwk.KeyType,
		Alg:              jwk.Alg,
		KeyID:            jwk.KeyID,
		Curve:            jwk.Curve,
		N:                jwk.N,
		E:                jwk.E,
		X:                jwk.X,
		Y:                jwk.Y,
		KeyOperations:    jwk.KeyOperations,
		ThumbprintSHA1:   jwk.ThumbprintSHA1,
		ThumbprintSHA256: jwk.ThumbprintSHA256,
		CertificateChain: jwk.CertificateChain,
	}
}

func (jwk *JsonWebKey) PublicKey() crypto.PublicKey {
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

func (jwk *JsonWebKey) rsaPublicKey() *rsa.PublicKey {
	if (jwk.cachedPublicKey) != nil {
		return jwk.cachedPublicKey.(*rsa.PublicKey)
	}
	jwk.cachedPublicKey = &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(jwk.N),
		E: int(big.NewInt(0).SetBytes(jwk.E).Int64()),
	}
	return jwk.cachedPublicKey.(*rsa.PublicKey)
}

func (jwk *JsonWebKey) ecdsaPublicKey() *ecdsa.PublicKey {
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
func (jwk *JsonWebKey) PrivateKey() crypto.PrivateKey {
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
				Dp:   big.NewInt(0).SetBytes(jwk.Dp),
				Dq:   big.NewInt(0).SetBytes(jwk.Dq),
				Qinv: big.NewInt(0).SetBytes(jwk.Qinv),
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

func (jwk *JsonWebKey) SetPublicKey(publicKey crypto.PublicKey) error {
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		jwk.KeyType = KeyTypeRSA
		jwk.N = publicKey.N.Bytes()
		jwk.E = big.NewInt(int64(publicKey.E)).Bytes()
		jwk.X = nil
		jwk.Y = nil
	case *ecdsa.PublicKey:
		jwk.KeyType = KeyTypeEC
		switch publicKey.Curve {
		case elliptic.P256():
			jwk.Curve = CurveNameP256
		case elliptic.P384():
			jwk.Curve = CurveNameP384
		case elliptic.P521():
			jwk.Curve = CurveNameP521
		default:
			return errInvalidCurve
		}
		jwk.X = publicKey.X.Bytes()
		jwk.Y = publicKey.Y.Bytes()
		jwk.N = nil
		jwk.E = nil
	default:
		return ErrInvalidKeyType
	}
	jwk.cachedPublicKey = publicKey
	jwk.cachedPrivateKey = nil
	jwk.D = nil
	jwk.P = nil
	jwk.Q = nil
	jwk.Dp = nil
	jwk.Dq = nil
	jwk.Qinv = nil
	return nil
}

func SanitizeKeyOperations(keyOps []JsonWebKeyOperation) []JsonWebKeyOperation {
	if keyOps == nil {
		return nil
	}
	seen := make(map[JsonWebKeyOperation]bool)
	result := make([]JsonWebKeyOperation, 0, len(keyOps))
	for _, op := range keyOps {
		if _, ok := seen[op]; ok {
			continue
		}
		seen[op] = true
		result = append(result, op)
	}
	return result
}

func NewJsonWebKeyFromPublicKey(publicKey crypto.PublicKey) (*JsonWebKey, error) {
	jwk := &JsonWebKey{}
	err := jwk.SetPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return jwk, nil
}
