package cert

import (
	"crypto"
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type azKeysSelfSignerProvider struct {
	certDoc *CertDoc

	cert           *x509.Certificate
	certSpec       CertJwkSpec
	signer         *keyVaultSigner
	keyCreated     *azkeys.ID
	keepKeyVersion bool
	eCtx           common.ElevatedContext
}

// X509SigningAlg implements SignerProvider.
func (p *azKeysSelfSignerProvider) X509SigningAlg() x509.SignatureAlgorithm {
	return p.certSpec.Alg.ToX509SignatureAlgorithm()
}

// Locator implements SignerProvider.
func (p *azKeysSelfSignerProvider) Locator() shared.ResourceLocator {
	return p.certDoc.GetLocator()
}

// Load implements CertificateRequestProvider.
func (p *azKeysSelfSignerProvider) Load(c common.ElevatedContext) (certTemplate *x509.Certificate, publicKey any, publicKeySpec *CertJwkSpec, err error) {

	bad := func(e error) (*x509.Certificate, any, *CertJwkSpec, error) {
		return nil, nil, nil, e
	}

	p.cert, err = p.certDoc.createX509Certificate()
	if err != nil {
		return bad(err)
	}
	p.certSpec = p.certDoc.CertSpec

	p.eCtx = c
	p.signer, p.keyCreated, err = createAzKeysSigner(p.eCtx, &p.certSpec,
		*p.certDoc.KeyStorePath,
		&azkeys.KeyAttributes{
			Expires:    p.certDoc.NotAfter.TimePtr(),
			Exportable: &p.certSpec.keyExportable,
		})
	if err != nil {
		return bad(err)
	}

	return p.cert, p.signer.publicKey, &p.certSpec, err
}

// LocatorString implements SignerProvider.
func (p *azKeysSelfSignerProvider) GetIssuerCertStorePath() string {
	return "@self"
}

// CollectCertificateChain implements CertificateRequestProvider.
func (p *azKeysSelfSignerProvider) CollectCertificateChain([][]byte, *CertJwkSpec) error {
	p.keepKeyVersion = true
	return nil
}

// Certificate implements SignerProvider.
func (p *azKeysSelfSignerProvider) Certificate() *x509.Certificate {
	return p.cert
}

// CertificatesInChain implements SignerProvider.
func (*azKeysSelfSignerProvider) CertificatesInChain() [][]byte {
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
		client := p.signer.keysClient
		_, err := client.UpdateKey(p.eCtx, p.keyCreated.Name(), p.keyCreated.Version(), azkeys.UpdateKeyParameters{
			KeyAttributes: &azkeys.KeyAttributes{
				Enabled: utils.ToPtr(false),
			},
		}, nil)
		if err != nil {
			log.Error().Err(err).Msgf("%s:failed to disable key version: %s", err, *p.keyCreated)
		}
	}
}

// Signer implements SignerProvider.
func (p *azKeysSelfSignerProvider) LoadSigner(common.ElevatedContext) (crypto.Signer, error) {
	return p.signer, nil
}

func createAzKeysSigner(c common.ElevatedContext, ioCertJwkSpec *CertJwkSpec, keyName string, keyAttributes *azkeys.KeyAttributes) (*keyVaultSigner, *azkeys.ID, error) {
	var keyCreated *azkeys.ID
	bad := func(e error) (*keyVaultSigner, *azkeys.ID, error) {
		return nil, keyCreated, e
	}

	client := common.GetAdminServerClientProvider(c).AzKeysClient()
	var signingAlg azkeys.SignatureAlgorithm
	params := azkeys.CreateKeyParameters{}
	switch ioCertJwkSpec.Kty {
	case models.KeyTypeRSA:
		params.Kty = utils.ToPtr(azkeys.KeyTypeRSA)
		params.KeySize = ioCertJwkSpec.KeySize
		signingAlg = azkeys.SignatureAlgorithmRS384
		switch ioCertJwkSpec.Alg {
		case shared.AlgRS256:
			signingAlg = azkeys.SignatureAlgorithmRS256
		case shared.AlgRS512:
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
