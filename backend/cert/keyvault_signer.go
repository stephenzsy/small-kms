package cert

import (
	"context"
	"crypto"
	"encoding/asn1"
	"io"
	"math/big"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/common"
)

// should implement crypto.Signer
type keyvaultSigner struct {
	ctx        context.Context
	keysClient *azkeys.Client
	jwk        *azkeys.JSONWebKey
	publicKey  crypto.PublicKey
	sigAlg     azkeys.SignatureAlgorithm
}

func newKeyVaultSignerWithExistingPublicKey(c context.Context, keyID *azkeys.ID, publicKey crypto.PublicKey, sigAlg azkeys.SignatureAlgorithm) *keyvaultSigner {
	return &keyvaultSigner{
		ctx:        c,
		keysClient: common.GetAdminServerClientProvider(c).AzKeysClient(),
		jwk: &azkeys.JSONWebKey{
			KID: keyID,
		},
		publicKey: publicKey,
		sigAlg:    sigAlg,
	}
}

type esSignature struct {
	R *big.Int
	S *big.Int
}

func (s *keyvaultSigner) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
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
