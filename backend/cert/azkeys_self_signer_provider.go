package cert

import (
	"crypto"
	"crypto/x509"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type azKeysSelfSignerProvider struct {
	certDoc *CertDoc

	cert           *x509.Certificate
	certSpec       CertJwkSpec
	signer         *keyVaultSigner
	keyCreated     *azkeys.ID
	keepKeyVersion bool
	ctx            common.ServiceContext
}

// Locator implements SignerProvider.
func (p *azKeysSelfSignerProvider) Locator() common.Locator[models.NamespaceKind, models.ResourceKind] {
	return p.certDoc.GetLocator()
}

// Load implements CertificateRequestProvider.
func (p *azKeysSelfSignerProvider) Load(c common.ServiceContext) (certTemplate *x509.Certificate, publicKey any, publicKeySpec *CertJwkSpec, err error) {

	bad := func(e error) (*x509.Certificate, any, *CertJwkSpec, error) {
		return nil, nil, nil, e
	}

	p.ctx = c
	p.cert, err = p.certDoc.createX509Certificate()
	if err != nil {
		return bad(err)
	}
	p.certSpec = p.certDoc.CertSpec

	p.signer, p.keyCreated, err = createAzKeysSigner(c, &p.certSpec,
		*p.certDoc.KeyStorePath,
		&azkeys.KeyAttributes{
			Expires:    p.certDoc.NotAfter.TimePtr(),
			Exportable: &p.certSpec.keyExportable,
		})
	if err != nil {
		return bad(err)
	}

	switch p.certSpec.Alg {
	case models.AlgRS256:
		p.cert.SignatureAlgorithm = x509.SHA256WithRSA
	case models.AlgRS384:
		p.cert.SignatureAlgorithm = x509.SHA384WithRSA
	case models.AlgRS512:
		p.cert.SignatureAlgorithm = x509.SHA512WithRSA
	case models.AlgES256:
		p.cert.SignatureAlgorithm = x509.ECDSAWithSHA256
	case models.AlgES384:
		p.cert.SignatureAlgorithm = x509.ECDSAWithSHA384
	default:
		return bad(fmt.Errorf("%w:unsupported cert signature algorithm:%s", common.ErrStatusBadRequest, p.certSpec.Alg))
	}

	return p.cert, p.signer.publicKey, &p.certSpec, err
}

// LocatorString implements SignerProvider.
func (p *azKeysSelfSignerProvider) GetIssuerCertStorePath() string {
	return "@self"
}

// CollectCertificateChain implements CertificateRequestProvider.
func (p *azKeysSelfSignerProvider) CollectCertificateChain([][]byte) error {
	p.keepKeyVersion = true
	return nil
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
	if !p.keepKeyVersion && p.keyCreated != nil {
		// disable key version
		client := common.GetClientProvider(p.ctx).AzKeysClient()
		_, err := client.UpdateKey(p.ctx, p.keyCreated.Name(), p.keyCreated.Version(), azkeys.UpdateKeyParameters{
			KeyAttributes: &azkeys.KeyAttributes{
				Enabled: utils.ToPtr(false),
			},
		}, nil)
		if err != nil {
			log.Error().Err(err).Msgf("failed to disable key version: %s", *p.keyCreated)
		}
	}
}

// Signer implements SignerProvider.
func (p *azKeysSelfSignerProvider) LoadSigner(common.ServiceContext) (crypto.Signer, error) {
	return p.signer, nil
}

func createAzKeysSigner(c common.ServiceContext, ioCertJwkSpec *CertJwkSpec, keyName string, keyAttributes *azkeys.KeyAttributes) (*keyVaultSigner, *azkeys.ID, error) {
	var keyCreated *azkeys.ID
	bad := func(e error) (*keyVaultSigner, *azkeys.ID, error) {
		return nil, keyCreated, e
	}

	client := common.GetClientProvider(c).AzKeysClient()
	var signingAlg azkeys.SignatureAlgorithm
	params := azkeys.CreateKeyParameters{}
	switch ioCertJwkSpec.Kty {
	case models.KeyTypeRSA:
		params.Kty = utils.ToPtr(azkeys.KeyTypeRSA)
		params.KeySize = ioCertJwkSpec.KeySize
		signingAlg = azkeys.SignatureAlgorithmRS384
		switch ioCertJwkSpec.Alg {
		case models.AlgRS256:
			signingAlg = azkeys.SignatureAlgorithmRS256
		case models.AlgRS512:
			signingAlg = azkeys.SignatureAlgorithmRS512
		}
	case models.KeyTypeEC:
		params.Kty = utils.ToPtr(azkeys.KeyTypeEC)
		params.Curve = utils.ToPtr(azkeys.CurveNameP384)
		signingAlg = azkeys.SignatureAlgorithmES384
		switch *ioCertJwkSpec.Crv {
		case models.CurveNameP256:
			params.Curve = utils.ToPtr(azkeys.CurveNameP256)
			signingAlg = azkeys.SignatureAlgorithmES256
		}
	}

	params.KeyAttributes = keyAttributes
	params.KeyOps = []*azkeys.KeyOperation{
		utils.ToPtr(azkeys.KeyOperationSign),
		utils.ToPtr(azkeys.KeyOperationVerify),
	}

	keyResp, err := client.CreateKey(c, keyName, params, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create key: %s", keyName)
		return bad(err)
	}
	jwk := keyResp.Key
	keyCreated = jwk.KID
	log.Info().Msgf("key created: %s", *keyCreated)

	signer, err := newKeyVaultSigner(c, client, jwk, signingAlg)
	if err != nil {
		return bad(err)
	}

	// update certJwkSpec
	ioCertJwkSpec.KID = string(*jwk.KID)
	return signer, keyCreated, nil
}

var _ SignerProvider = (*azKeysSelfSignerProvider)(nil)
var _ CertificateRequestProvider = (*azKeysSelfSignerProvider)(nil)

func newAzKeysSelfSignerProvider(certDoc *CertDoc) *azKeysSelfSignerProvider {
	return &azKeysSelfSignerProvider{
		certDoc: certDoc,
	}
}
