package cert

import (
	"bytes"
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateRequestProvider interface {
	Load(common.ServiceContext) (certTemplate *x509.Certificate, publicKey any, publicKeySpec *CertJwkSpec, err error)
	Close()
	CollectCertificateChain([][]byte) error
}

type SignerProvider interface {
	// this call also populate other fields in the signer provider
	LoadSigner(common.ServiceContext) (crypto.Signer, error)
	Certificate() *x509.Certificate
	Locator() models.ResourceLocator
	GetIssuerCertStorePath() string
	CertificateChainPEM() []byte
	ExtraCertificatesInChain() [][]byte
}

type StorageProvider interface {
	StoreCertificateChainPEM(c common.ServiceContext, pemBlob []byte, x5t []byte,
		issuerLocatorStr string) (string, error)
}

func signCertificate(c common.ServiceContext,
	csrProvider CertificateRequestProvider,
	signerProvider SignerProvider,
	storageProvider StorageProvider) (*CertDocSigningPatch, error) {

	// load certificate public key first, in case of our implementation of self signer requires key created before signing
	defer csrProvider.Close()
	certTemplate, publicKey, certJwkSpec, err := csrProvider.Load(c)

	if err != nil {
		return nil, err
	}

	signer, err := signerProvider.LoadSigner(c)
	if err != nil {
		return nil, err
	}

	certCreated, err := x509.CreateCertificate(nil,
		certTemplate,
		signerProvider.Certificate(),
		publicKey,
		signer)
	if err != nil {
		return nil, err
	}
	fullChain := append([][]byte{certCreated}, signerProvider.ExtraCertificatesInChain()...)
	err = csrProvider.CollectCertificateChain(fullChain)
	if err != nil {
		return nil, err
	}

	pemBuf := bytes.Buffer{}
	err = pem.Encode(&pemBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certCreated})
	if err != nil {
		return nil, err
	}
	pemBuf.Write(signerProvider.CertificateChainPEM())
	x5t := sha1.Sum(certCreated)
	x5tS256 := sha256.Sum256(certCreated)
	certJwkSpec.X5t = x5t[:]
	certJwkSpec.X5tS256 = x5tS256[:]
	blobKey, err := storageProvider.StoreCertificateChainPEM(c, pemBuf.Bytes(), x5t[:], signerProvider.GetIssuerCertStorePath())
	if err != nil {
		return nil, err
	}
	return &CertDocSigningPatch{
		CertSpec:      *certJwkSpec,
		Thumbprint:    x5t[:],
		CertStorePath: blobKey,
		Issuer:        signerProvider.Locator(),
	}, nil
}

type azBlobStorageProvider struct {
	blobKey string
}

func (p *azBlobStorageProvider) StoreCertificateChainPEM(c common.ServiceContext, pem, x5t []byte, issuerLocaterStr string) (string, error) {
	if p.blobKey == "" {
		return "", fmt.Errorf("%w:empty blob name", common.ErrStatusBadRequest)
	}
	blockBlobClient := common.GetClientProvider(c).AzBlobContainerClient().NewBlockBlobClient(p.blobKey)
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
