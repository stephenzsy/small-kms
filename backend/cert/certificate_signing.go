package cert

import (
	"bytes"
	"context"
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
	"github.com/stephenzsy/small-kms/backend/shared"
)

type SelfSignedCertificateProvider interface {
	CreateSelfSignedCertificate(context.Context) ([]byte, *CertJwkSpec, error)
	Close(context.Context)
	KeepCertificate()
	Locator() shared.ResourceLocator
}

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

func getSelfSignedCertificate(c context.Context,
	certProvider SelfSignedCertificateProvider,
	storageProvider StorageProvider) (*CertDocSigningPatch, error) {
	defer certProvider.Close(c)
	certCreated, certJwkSpec, err := certProvider.CreateSelfSignedCertificate(c)

	if err != nil {
		return nil, err
	}

	pemBuf := bytes.Buffer{}
	err = pem.Encode(&pemBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certCreated})
	if err != nil {
		return nil, err
	}
	x5t := sha1.Sum(certCreated)
	x5tS256 := sha256.Sum256(certCreated)
	certJwkSpec.X5t = x5t[:]
	certJwkSpec.X5tS256 = x5tS256[:]
	blobKey, err := storageProvider.StoreCertificateChainPEM(c, pemBuf.Bytes(), x5t[:], "@self")
	if err != nil {
		return nil, err
	}
	return &CertDocSigningPatch{
		CertSpec:      *certJwkSpec,
		Thumbprint:    x5t[:],
		CertStorePath: blobKey,
		Issuer:        certProvider.Locator(),
	}, nil
}

func signCertificate(c context.Context,
	csrProvider CertificateRequestProvider,
	signerProvider SignerProvider,
	storageProvider StorageProvider) (*CertDocSigningPatch, error) {

	// load certificate public key first, in case of our implementation of self signer requires key created before signing
	defer csrProvider.Close(c)
	certTemplate, publicKey, certJwkSpec, err := csrProvider.Load(c)

	if err != nil {
		return nil, err
	}

	signer, err := signerProvider.LoadSigner(c)
	if err != nil {
		return nil, err
	}
	certTemplate.SignatureAlgorithm = signerProvider.X509SigningAlg()

	certCreated, err := x509.CreateCertificate(nil,
		certTemplate,
		signerProvider.Certificate(),
		publicKey,
		signer)
	if err != nil {
		return nil, err
	}
	fullChain := append([][]byte{certCreated}, signerProvider.CertificatesInChain()...)
	err = csrProvider.CollectCertificateChain(c, fullChain, certJwkSpec)
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
