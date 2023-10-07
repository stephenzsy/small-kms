package cert

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"

	"github.com/stephenzsy/small-kms/backend/common"
)

type CertificateRequestProvider interface {
	PublicKey() any
	Close()
	CollectCertificateChain([][]byte) error
}

type SignerProvider interface {
	Certificate() *x509.Certificate
	Signer() crypto.Signer
	Close()
	CertificateChainPEM() []byte
	CertificateChain() [][]byte
}

type CertificateFieldsProvider interface {
	PopulateX509(cert *x509.Certificate) error
}

type StorageProvider interface {
	StoreCertificateChainPEM([]byte) error
}

func signCertificate(c common.ServiceContext,
	csrProvider CertificateRequestProvider,
	signerProvider SignerProvider,
	certificateFieldsProvider CertificateFieldsProvider,
	storageProvider StorageProvider) ([]byte, error) {
	certTemplate := x509.Certificate{}
	err := certificateFieldsProvider.PopulateX509(&certTemplate)
	if err != nil {
		return nil, err
	}

	defer csrProvider.Close()
	defer signerProvider.Close()

	publicKey := csrProvider.PublicKey()
	certCreated, err := x509.CreateCertificate(nil, &certTemplate, signerProvider.Certificate(), publicKey, signerProvider.Signer())
	if err != nil {
		return nil, err
	}
	fullChain := append([][]byte{certCreated}, signerProvider.CertificateChain()...)
	csrProvider.CollectCertificateChain(fullChain)
	pemBuf := bytes.Buffer{}
	err = pem.Encode(&pemBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certCreated})
	if err != nil {
		return certCreated, err
	}
	pemBuf.Write(signerProvider.CertificateChainPEM())
	err = storageProvider.StoreCertificateChainPEM(pemBuf.Bytes())
	if err != nil {
		return certCreated, err
	}
	return certCreated, nil
}
