package cert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type azCertsCsrProvider struct {
	certDoc *CertDoc

	client                  *azcertificates.Client
	cert                    *x509.Certificate
	csr                     *x509.CertificateRequest
	certOperationInProgress string
	selfSignedCertId        *azcertificates.ID
	eCtx                    common.ElevatedContext
	selfSigned              bool
}

// CreateSelfSignedCertificate implements SelfSignedCertificateProvider.
func (p *azCertsCsrProvider) CreateSelfSignedCertificate(c common.ElevatedContext) ([]byte, *CertJwkSpec, error) {
	resp, err := p.createCert(c, true)
	if err != nil {
		return nil, nil, err
	}
	p.selfSignedCertId = resp.ID
	certResp, err := p.client.GetCertificate(p.eCtx, resp.ID.Name(), resp.ID.Version(), nil)
	if err != nil {
		return nil, nil, err
	}
	return certResp.CER, &p.certDoc.CertSpec, nil
}

// KeepCertificate implements SelfSignedCertificateProvider.
func (p *azCertsCsrProvider) KeepCertificate() {
	p.selfSignedCertId = nil
}

// Locator implements SelfSignedCertificateProvider.
func (p *azCertsCsrProvider) Locator() shared.ResourceLocator {
	return p.certDoc.GetLocator()
}

// Close implements CertificateRequestProvider.
func (p *azCertsCsrProvider) Close() {
	if p.client != nil {
		if p.selfSigned && p.selfSignedCertId != nil {
			p.client.UpdateCertificate(p.eCtx, p.selfSignedCertId.Name(), p.selfSignedCertId.Version(), azcertificates.UpdateCertificateParameters{
				CertificateAttributes: &azcertificates.CertificateAttributes{
					Enabled: utils.ToPtr(false),
				},
			}, nil)
		} else if !p.selfSigned && p.certOperationInProgress != "" {
			p.client.DeleteCertificate(p.eCtx, p.certOperationInProgress, nil)
		}
	}
}

// CollectCertificateChain implements CertificateRequestProvider.
func (p *azCertsCsrProvider) CollectCertificateChain(x5c [][]byte, ioCertSpec *CertJwkSpec) error {
	resp, err := p.client.MergeCertificate(p.eCtx, p.certOperationInProgress, azcertificates.MergeCertificateParameters{
		X509Certificates: x5c,
	}, nil)
	if err != nil {
		return err
	}
	p.certOperationInProgress = ""
	ioCertSpec.KID = string(*resp.KID)
	ioCertSpec.X5u = utils.ToPtr(string(*resp.ID))
	return nil
}

func (p *azCertsCsrProvider) createCert(c common.ElevatedContext, selfSigned bool) (resp azcertificates.CreateCertificateResponse, err error) {
	bad := func(e error) (azcertificates.CreateCertificateResponse, error) {
		return azcertificates.CreateCertificateResponse{}, e
	}

	p.cert, err = p.certDoc.createX509Certificate()
	if err != nil {
		return bad(err)
	}

	keyProperties := azcertificates.KeyProperties{
		Exportable: utils.ToPtr(p.certDoc.CertSpec.keyExportable),
	}

	switch p.certDoc.CertSpec.Kty {
	case shared.KeyTypeRSA:
		keyProperties.KeyType = utils.ToPtr(azcertificates.KeyTypeRSA)
		keyProperties.KeySize = p.certDoc.CertSpec.KeySize
	case shared.KeyTypeEC:
		keyProperties.KeyType = utils.ToPtr(azcertificates.KeyTypeEC)
		switch *p.certDoc.CertSpec.Crv {
		case shared.CurveNameP256:
			keyProperties.Curve = utils.ToPtr(azcertificates.CurveNameP256)
		case shared.CurveNameP384:
			keyProperties.Curve = utils.ToPtr(azcertificates.CurveNameP384)
		default:
			return bad(fmt.Errorf("unsupported curve: %s", *p.certDoc.CertSpec.Crv))
		}
	default:
		return bad(fmt.Errorf("unsupported key type: %s", p.certDoc.CertSpec.Kty))
	}

	csp := azcertificates.CreateCertificateParameters{
		CertificatePolicy: &azcertificates.CertificatePolicy{
			KeyProperties: &keyProperties,
			X509CertificateProperties: &azcertificates.X509CertificateProperties{
				Subject: utils.ToPtr(p.cert.Subject.String()),
			},
			SecretProperties: &azcertificates.SecretProperties{
				ContentType: to.Ptr("application/x-pem-file"),
			},
		},
	}
	if selfSigned {
		csp.CertificatePolicy.IssuerParameters = &azcertificates.IssuerParameters{
			Name: utils.ToPtr("Self"),
		}
	}

	p.client = common.GetAdminServerClientProvider(c).AzCertificatesClient()
	certName := *p.certDoc.KeyStorePath
	p.eCtx = c
	resp, err = p.client.CreateCertificate(p.eCtx, certName, csp, nil)
	return
}

// Load implements CertificateRequestProvider.
func (p *azCertsCsrProvider) Load(c common.ElevatedContext) (certTemplate *x509.Certificate, publicKey any, publicKeySpec *CertJwkSpec, err error) {
	bad := func(e error) (*x509.Certificate, any, *CertJwkSpec, error) {
		return nil, nil, nil, e
	}

	resp, err := p.createCert(c, false)
	if err != nil {
		return bad(err)
	}
	p.certOperationInProgress = resp.ID.Name()
	p.csr, err = x509.ParseCertificateRequest(resp.CSR)
	if err != nil {
		return bad(err)
	}
	p.cert.PublicKey = p.csr.PublicKey

	return p.cert, p.csr.PublicKey, &p.certDoc.CertSpec, nil
}

var _ CertificateRequestProvider = (*azCertsCsrProvider)(nil)
var _ SelfSignedCertificateProvider = (*azCertsCsrProvider)(nil)

func newAzCertsCsrProvider(certDoc *CertDoc, selfSigned bool) *azCertsCsrProvider {
	return &azCertsCsrProvider{
		certDoc:    certDoc,
		selfSigned: selfSigned,
	}
}

type azKeysExistingCertSigner struct {
	issuerCertDoc *CertDoc

	issuerCert       *x509.Certificate
	certChainPemBlob []byte
	restX5c          [][]byte
	keyVaultSigner   *keyVaultSigner
}

// X509SigningAlg implements SignerProvider.
func (p *azKeysExistingCertSigner) X509SigningAlg() x509.SignatureAlgorithm {
	return p.issuerCertDoc.CertSpec.Alg.ToX509SignatureAlgorithm()
}

// Certificate implements SignerProvider.
func (p *azKeysExistingCertSigner) Certificate() *x509.Certificate {
	return p.issuerCert
}

// CertificateChainPEM implements SignerProvider.
func (p *azKeysExistingCertSigner) CertificateChainPEM() []byte {
	return p.certChainPemBlob
}

// CertificatesInChain implements SignerProvider.
func (p *azKeysExistingCertSigner) CertificatesInChain() [][]byte {
	return p.restX5c
}

// GetIssuerCertStorePath implements SignerProvider.
func (p *azKeysExistingCertSigner) GetIssuerCertStorePath() string {
	return p.issuerCertDoc.CertStorePath
}

// LoadSigner implements SignerProvider.
func (p *azKeysExistingCertSigner) LoadSigner(c common.ElevatedContext) (signer crypto.Signer, err error) {
	kidStr := p.issuerCertDoc.CertSpec.KID
	if kidStr == "" {
		return nil, fmt.Errorf("empty key id from issuer")
	}
	p.certChainPemBlob, err = p.issuerCertDoc.fetchCertificatePEMBlob(c)
	if err != nil {
		return nil, err
	}

	p.restX5c = make([][]byte, 0, 2)
	for block, rest := pem.Decode(p.certChainPemBlob); block != nil; block, rest = pem.Decode(rest) {
		p.restX5c = append(p.restX5c, block.Bytes)
	}

	p.issuerCert, err = x509.ParseCertificate(p.restX5c[0])
	if err != nil {
		return nil, err
	}

	p.keyVaultSigner = &keyVaultSigner{
		ctx:        c,
		keysClient: common.GetAdminServerClientProvider(c).AzKeysClient(),
		jwk: &azkeys.JSONWebKey{
			KID: utils.ToPtr(azkeys.ID(p.issuerCertDoc.CertSpec.KID)),
		},
		publicKey: p.issuerCert.PublicKey,
		sigAlg:    p.issuerCertDoc.CertSpec.Alg.ToAzKeysSignatureAlgorithm(),
	}

	return p.keyVaultSigner, nil
}

// Locator implements SignerProvider.
func (p *azKeysExistingCertSigner) Locator() shared.ResourceLocator {
	return p.issuerCertDoc.GetLocator()
}

var _ SignerProvider = (*azKeysExistingCertSigner)(nil)

func newAzKeysExistingCertSigner(issuerCertDoc *CertDoc) *azKeysExistingCertSigner {
	return &azKeysExistingCertSigner{
		issuerCertDoc: issuerCertDoc,
	}
}
