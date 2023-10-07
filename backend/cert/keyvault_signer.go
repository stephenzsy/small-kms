package cert

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

// should implement crypto.Signer
type keyVaultSigner struct {
	ctx        context.Context
	keysClient *azkeys.Client
	jwk        *azkeys.JSONWebKey
	publicKey  crypto.PublicKey
	sigAlg     azkeys.SignatureAlgorithm
}

func newKeyVaultSigner(ctx context.Context, keysClient *azkeys.Client, jwk *azkeys.JSONWebKey, preferredRsaSigAlg azkeys.SignatureAlgorithm) (signer *keyVaultSigner, err error) {
	signer = &keyVaultSigner{
		ctx:        ctx,
		keysClient: keysClient,
		jwk:        jwk,
	}
	signer.publicKey, err = extractPublicKey(jwk)
	if err != nil {
		return nil, err
	}
	switch *jwk.Kty {
	case azkeys.KeyTypeRSA:
		signer.sigAlg = azkeys.SignatureAlgorithmRS384
		switch preferredRsaSigAlg {
		case azkeys.SignatureAlgorithmRS256:
			signer.sigAlg = azkeys.SignatureAlgorithmRS256
		case azkeys.SignatureAlgorithmRS512:
			signer.sigAlg = azkeys.SignatureAlgorithmRS512
		}
	case azkeys.KeyTypeEC:
		switch *jwk.Crv {
		case azkeys.CurveNameP256:
			signer.sigAlg = azkeys.SignatureAlgorithmES256
		case azkeys.CurveNameP384:
			signer.sigAlg = azkeys.SignatureAlgorithmES384
		default:
			return nil, fmt.Errorf("unsupported curve: %s", *jwk.Crv)
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %s", *jwk.Kty)
	}
	return signer, nil
}

func (s *keyVaultSigner) Public() crypto.PublicKey {
	return s.publicKey
}

type esSignature struct {
	R *big.Int
	S *big.Int
}

func (s *keyVaultSigner) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	resp, err := s.keysClient.Sign(s.ctx, s.jwk.KID.Name(), s.jwk.KID.Version(), azkeys.SignParameters{
		Algorithm: &s.sigAlg,
		Value:     digest,
	}, nil)
	if err != nil {
		return
	}
	signature = resp.Result
	n := 0
	switch s.sigAlg {
	case azkeys.SignatureAlgorithmES256:
		n = 32
	case azkeys.SignatureAlgorithmES384:
		n = 48
	}
	if n != 0 {
		sig := esSignature{
			R: new(big.Int).SetBytes(signature[:n]),
			S: new(big.Int).SetBytes(signature[n:]),
		}
		signature, err = asn1.Marshal(sig)
		if err != nil {
			return
		}
	}
	return
}

func extractPublicKey(key *azkeys.JSONWebKey) (crypto.PublicKey, error) {
	if *key.Kty == azkeys.KeyTypeRSA {
		k := &rsa.PublicKey{}

		// N = modulus
		if len(key.N) == 0 {
			return nil, errors.New("property N is empty")
		}
		k.N = &big.Int{}
		k.N = k.N.SetBytes(key.N)

		// e = public exponent
		if len(key.E) == 0 {
			return nil, errors.New("property e is empty")
		}
		k.E = int(big.NewInt(0).SetBytes(key.E).Uint64())
		return k, nil
	} else if *key.Kty == azkeys.KeyTypeEC {
		k := &ecdsa.PublicKey{}

		switch *key.Crv {
		case azkeys.CurveNameP256:
			k.Curve = elliptic.P256()
		case azkeys.CurveNameP384:
			k.Curve = elliptic.P384()
		default:
			return nil, fmt.Errorf("unsupported curve: %s", *key.Crv)
		}

		if len(key.X) == 0 {
			return nil, errors.New("property X is empty")
		}
		k.X = big.NewInt(0).SetBytes(key.X)

		if len(key.Y) == 0 {
			return nil, errors.New("property Y is empty")
		}
		k.Y = big.NewInt(0).SetBytes(key.Y)

		return k, nil
	}
	return nil, fmt.Errorf("unsupported key type: %s", *key.Kty)
}

var _ crypto.Signer = (*keyVaultSigner)(nil)
