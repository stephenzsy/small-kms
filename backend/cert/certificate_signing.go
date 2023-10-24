package cert

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type CertificateRequestProvider interface {
	Load(context.Context) (certTemplate *x509.Certificate, publicKey any, publicKeySpec *CertJwkSpec, err error)
	Close(context.Context)
	CollectCertificateChain(context.Context, [][]byte, *CertJwkSpec) error
}

type SignerProvider interface {
	// this call also populate other fields in the signer provider
	LoadSigner(context.Context) (crypto.Signer, error)
	Certificate() *x509.Certificate
	Locator() shared.ResourceLocator
	GetIssuerCertStorePath() string
	CertificateChainPEM() []byte
	CertificatesInChain() [][]byte
	X509SigningAlg() x509.SignatureAlgorithm
}

type StorageProvider interface {
	StoreCertificateChainPEM(c context.Context, pemBlob []byte, x5t []byte,
		issuerLocatorStr string) (string, error)
}

type azBlobStorageProvider struct {
	blobKey string
}

func (p *azBlobStorageProvider) StoreCertificateChainPEM(c context.Context, pem, x5t []byte, issuerLocaterStr string) (string, error) {
	if p.blobKey == "" {
		return "", fmt.Errorf("%w:empty blob name", common.ErrStatusBadRequest)
	}
	blockBlobClient := common.GetAdminServerClientProvider(c).CertsAzBlobContainerClient().NewBlockBlobClient(p.blobKey)
	_, err := blockBlobClient.UploadBuffer(c, pem, &blockblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr("application/x-pem-file"),
		},
		Metadata: map[string]*string{
			"issuer": &issuerLocaterStr,
			"x5t":    to.Ptr(hex.EncodeToString(x5t)),
		},
	})
	return p.blobKey, err
}

var _ StorageProvider = (*azBlobStorageProvider)(nil)
