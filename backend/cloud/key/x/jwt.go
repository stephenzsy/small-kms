package cloudkeyx

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"

	"github.com/golang-jwt/jwt/v5"
	i "github.com/stephenzsy/small-kms/backend/cloud/key"
)

type cloudKeySigningMethod struct {
	alg i.JsonWebSignatureAlgorithm
}

// Alg implements jwt.SigningMethod.
func (m *cloudKeySigningMethod) Alg() string {
	return string(m.alg)
}

// Sign implements jwt.SigningMethod.
func (m *cloudKeySigningMethod) Sign(signingString string, key interface{}) ([]byte, error) {
	cloudKey, ok := key.(i.CloudSignatureKey)
	if !ok {
		return nil, fmt.Errorf("%w: %T", jwt.ErrInvalidKeyType, key)
	}
	hashFn := m.alg.HashFunc()
	if !hashFn.Available() {
		return nil, jwt.ErrHashUnavailable
	}

	digester := hashFn.New()
	digester.Write([]byte(signingString))
	return cloudKey.Sign(rand.Reader, digester.Sum(nil), m.alg)
}

func verifyEcdsaJWTSignature(pubKey *ecdsa.PublicKey, digest, signature []byte, keyByteSize int) error {

	r := big.NewInt(0).SetBytes(signature[:keyByteSize])
	s := big.NewInt(0).SetBytes(signature[keyByteSize:])

	if !ecdsa.Verify(pubKey,
		digest,
		r, s) {
		return jwt.ErrECDSAVerification
	}
	return nil
}

// This method does not support symmetric signing. use builtin JWT signing method instead.
func (m *cloudKeySigningMethod) Verify(signingString string, signature []byte, key interface{}) error {
	hashFn := m.alg.HashFunc()
	if !hashFn.Available() {
		return jwt.ErrHashUnavailable
	}
	digester := hashFn.New()
	digester.Write([]byte(signingString))
	digest := digester.Sum(nil)
	switch m.alg {
	case i.SignatureAlgorithmRS256, i.SignatureAlgorithmRS384, i.SignatureAlgorithmRS512:
		if pubKey, ok := key.(*rsa.PublicKey); ok {
			return rsa.VerifyPKCS1v15(pubKey, hashFn, digest, signature)
		}
	case i.SignatureAlgorithmPS256, i.SignatureAlgorithmPS384, i.SignatureAlgorithmPS512:
		if pubKey, ok := key.(*rsa.PublicKey); ok {
			return rsa.VerifyPSS(pubKey, hashFn, digest, signature, nil)
		}
	case i.SignatureAlgorithmES256:
		if pubKey, ok := key.(*ecdsa.PublicKey); ok {
			return verifyEcdsaJWTSignature(pubKey, digest, signature, 32)
		}
	case i.SignatureAlgorithmES384:
		if pubKey, ok := key.(*ecdsa.PublicKey); ok {
			return verifyEcdsaJWTSignature(pubKey, digest, signature, 48)
		}
	case i.SignatureAlgorithmES512:
		if pubKey, ok := key.(*ecdsa.PublicKey); ok {
			return verifyEcdsaJWTSignature(pubKey, digest, signature, 66)
		}
	}
	return jwt.ErrInvalidKeyType
}

var _ jwt.SigningMethod = (*cloudKeySigningMethod)(nil)

func NewJWTSigningMethod(alg i.JsonWebSignatureAlgorithm) jwt.SigningMethod {
	return &cloudKeySigningMethod{alg: alg}
}
