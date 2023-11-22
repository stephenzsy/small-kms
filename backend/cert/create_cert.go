package cert

import (
	"context"
	"crypto"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
)

type existingPublicKeyProvider struct {
	crypto.PublicKey
}

// Cleanup implements kv.AzCertCSRProvider.
func (*existingPublicKeyProvider) Cleanup(context.Context) {
	// do nothing
}

// CollectCerts implements kv.AzCertCSRProvider.
func (*existingPublicKeyProvider) CollectCerts(context.Context, [][]byte) (*azcertificates.MergeCertificateResponse, error) {
	return nil, nil
}

// GetCSRPublicKey implements kv.AzCertCSRProvider.
func (p *existingPublicKeyProvider) GetCSRPublicKey(context.Context) (crypto.PublicKey, error) {
	return p.PublicKey, nil
}

var _ kv.AzCertCSRProvider = (*existingPublicKeyProvider)(nil)
