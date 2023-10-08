package cert

import (
	"crypto"
	"crypto/x509"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type azKeysSelfSignerProvider struct {
	cert         *x509.Certificate
	keyStorePath string
	signer       *keyVaultSigner
	notAfter     time.Time
	certSpec     CertJwtSpec
	locator      models.ResourceLocator
}

// KeySpec implements CertificateRequestProvider.
func (p *azKeysSelfSignerProvider) KeySpec() CertJwtSpec {
	return p.certSpec
}

// LocatorString implements SignerProvider.
func (p *azKeysSelfSignerProvider) Locator() ResourceLocator {
	return p.locator
}

// CollectCertificateChain implements CertificateRequestProvider.
func (*azKeysSelfSignerProvider) CollectCertificateChain([][]byte) error {
	return nil
}

// PublicKey implements CertificateRequestProvider.
func (p *azKeysSelfSignerProvider) PublicKey() any {
	if p.signer != nil {
		return p.signer.Public()
	}
	return nil
}

// setCertificateTemplate implements SignerProvider.
func (p *azKeysSelfSignerProvider) setCertificateTemplate(cert *x509.Certificate) {
	p.cert = cert
}

// Certificate implements SignerProvider.
func (p *azKeysSelfSignerProvider) Certificate() *x509.Certificate {
	return p.cert
}

// ExtraCertificatesInChain implements SignerProvider.
func (*azKeysSelfSignerProvider) ExtraCertificatesInChain() [][]byte {
	return nil
}

// CertificateChainPEM implements SignerProvider.
func (*azKeysSelfSignerProvider) CertificateChainPEM() []byte {
	return nil
}

// Close implements SignerProvider.
func (p *azKeysSelfSignerProvider) Close() {
	p.signer = nil
}

// Signer implements SignerProvider.
func (p *azKeysSelfSignerProvider) GetSigner(c common.ServiceContext) (crypto.Signer, error) {
	if p.signer != nil {
		return p.signer, nil
	}

	client := common.GetClientProvider(c).AzKeysClient()
	var signingAlg azkeys.SignatureAlgorithm
	params := azkeys.CreateKeyParameters{}
	switch p.certSpec.Kty {
	case models.KeyTypeRSA:
		params.Kty = utils.ToPtr(azkeys.KeyTypeRSA)
		params.KeySize = p.certSpec.KeySize
		signingAlg = azkeys.SignatureAlgorithmRS384
		switch p.certSpec.Alg {
		case models.AlgRS256:
			signingAlg = azkeys.SignatureAlgorithmRS256
		case models.AlgRS512:
			signingAlg = azkeys.SignatureAlgorithmRS512
		}
	case models.KeyTypeEC:
		params.Kty = utils.ToPtr(azkeys.KeyTypeEC)
		params.Curve = utils.ToPtr(azkeys.CurveNameP384)
		signingAlg = azkeys.SignatureAlgorithmES384
		switch *p.certSpec.Crv {
		case models.CurveNameP256:
			params.Curve = utils.ToPtr(azkeys.CurveNameP256)
			signingAlg = azkeys.SignatureAlgorithmES256
		}
	}

	params.KeyAttributes = &azkeys.KeyAttributes{
		Expires:    &p.notAfter,
		Exportable: &p.certSpec.keyExportable,
	}

	keyResp, err := client.CreateKey(c, p.keyStorePath, params, nil)
	if err != nil {
		return nil, err
	}
	jwk := keyResp.Key

	return newKeyVaultSigner(c, client, jwk, signingAlg)
}

var _ SignerProvider = (*azKeysSelfSignerProvider)(nil)
var _ CertificateRequestProvider = (*azKeysSelfSignerProvider)(nil)
