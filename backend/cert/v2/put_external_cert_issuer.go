package cert

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"golang.org/x/crypto/acme"
)

func (*CertServer) PutExternalCertificateIssuer(ec echo.Context, namespaceId string, issuerID string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	req := &certmodels.CertificateExternalIssuerFields{}
	if err := c.Bind(req); err != nil {
		return err
	}

	acmeReq := req.Acme
	if acmeReq == nil {
		return base.ErrResponseStatusBadRequest
	}

	logger := log.Ctx(c)

	// load cloudKey
	keyDoc, err := key.GetKeyInternal(c, models.NamespaceProviderExternalCA, namespaceId, acmeReq.AccountKeyID)
	if err != nil {
		return err
	}
	ck := cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, kv.GetAzKeyVaultService(c).AzKeysClient(), keyDoc.KeyID, cloudkey.SignatureAlgorithmES384, true, keyDoc.PublicKey())

	acmeClient := acme.Client{
		Key:          ck,
		DirectoryURL: acmeReq.DirectoryURL,
	}
	account, err := acmeClient.GetReg(c, "")
	// account, err := acmeClient.Register(c, &acme.Account{
	// 	Contact: utils.MapSlice(acmeReq.Contacts, func(contact string) string {
	// 		return "mailto:" + contact
	// 	}),
	// }, func(tosURL string) bool {
	// 	logger.Info().Str("tosURL", tosURL).Msg("tosURL")
	// 	return true
	// })
	if err != nil {
		logger.Error().Err(err).Msg("acmeClient.Register")
		return err
	}

	issuerDoc := &CertIssuerDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderExternalCA,
				NamespaceID:       namespaceId,
				ResourceProvider:  models.ResourceProviderCertExternalIssuer,
			},
			ID: issuerID,
		},
		DisplayName: issuerID,
		ACME: &CertIssuerDocACME{
			AccountURI:     account.URI,
			AccountKeyID:   acmeReq.AccountKeyID,
			DirectoryURL:   acmeReq.DirectoryURL,
			AccountContact: account.Contact,
			AccountStatus:  account.Status,
		},
	}

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, issuerDoc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, issuerDoc.ToModel())
}
