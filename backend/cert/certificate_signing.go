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
	PublicKey() any
	Close()
	CollectCertificateChain([][]byte) error
	KeySpec() CertJwtSpec
}

type SignerProvider interface {
	Certificate() *x509.Certificate
	GetSigner(common.ServiceContext) (crypto.Signer, error)
	Locator() ResourceLocator
	Close()
	CertificateChainPEM() []byte
	ExtraCertificatesInChain() [][]byte
	// used only for self signing
	setCertificateTemplate(*x509.Certificate)
}

type CertificateFieldsProvider interface {
	PopulateX509(cert *x509.Certificate) error
}

type StorageProvider interface {
	StoreCertificateChainPEM(c common.ServiceContext, pemBlob []byte, x5t []byte,
		issuerLocatorStr string) (string, error)
}

func signCertificate(c common.ServiceContext,
	csrProvider CertificateRequestProvider,
	signerProvider SignerProvider,
	certificateFieldsProvider CertificateFieldsProvider,
	storageProvider StorageProvider) (*CertDocSigningPatch, error) {
	certTemplate := x509.Certificate{}
	err := certificateFieldsProvider.PopulateX509(&certTemplate)
	if err != nil {
		return nil, err
	}

	defer signerProvider.Close()
	signerProvider.setCertificateTemplate(&certTemplate)
	signer, err := signerProvider.GetSigner(c)
	if err != nil {
		return nil, err
	}

	defer csrProvider.Close()
	publicKey := csrProvider.PublicKey()
	certSpec := csrProvider.KeySpec()
	switch certSpec.Alg {
	case models.AlgRS256:
		certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
	case models.AlgRS384:
		certTemplate.SignatureAlgorithm = x509.SHA384WithRSA
	case models.AlgRS512:
		certTemplate.SignatureAlgorithm = x509.SHA512WithRSA
	case models.AlgES256:
		certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
	case models.AlgES384:
		certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA384
	default:
		return nil, fmt.Errorf("%w:unsupported cert signature algorithm:%s", common.ErrStatusBadRequest, certSpec.Alg)
	}

	certCreated, err := x509.CreateCertificate(nil,
		&certTemplate,
		signerProvider.Certificate(),
		publicKey,
		signer)
	if err != nil {
		return nil, err
	}
	fullChain := append([][]byte{certCreated}, signerProvider.ExtraCertificatesInChain()...)
	csrProvider.CollectCertificateChain(fullChain)
	pemBuf := bytes.Buffer{}
	err = pem.Encode(&pemBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certCreated})
	if err != nil {
		return nil, err
	}
	pemBuf.Write(signerProvider.CertificateChainPEM())
	x5t := sha1.Sum(certCreated)
	x5ts256 := sha256.Sum256(certCreated)
	blobKey, err := storageProvider.StoreCertificateChainPEM(c, pemBuf.Bytes(), x5t[:], signerProvider.Locator().String())
	if err != nil {
		return nil, err
	}
	return &CertDocSigningPatch{
		CertSpec: CertJwtSpec{
			CertKeySpec: certSpec.CertKeySpec,
			X5t:         x5t[:],
			X5tS256:     x5ts256[:],
		},
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
			"issuer_id": to.Ptr(issuerLocaterStr),
			"x5t":       to.Ptr(hex.EncodeToString(x5t)),
		},
	})
	return p.blobKey, err
}

var _ StorageProvider = (*azBlobStorageProvider)(nil)
