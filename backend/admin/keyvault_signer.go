package admin

import (
	"context"
	"crypto"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

type keyVaultSigner struct {
	ctx        context.Context
	keysClient *azkeys.Client
	kid        azkeys.ID
	publicKey  crypto.PublicKey
}

func (s *keyVaultSigner) Public() crypto.PublicKey {
	return s.publicKey
}

func (s *keyVaultSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	resp, err := s.keysClient.Sign(s.ctx, s.kid.Name(), s.kid.Version(), azkeys.SignParameters{
		Algorithm: to.Ptr(azkeys.SignatureAlgorithmRS384),
		Value:     digest,
	}, nil)
	if err != nil {
		return
	}
	signature = resp.Result
	return
}
