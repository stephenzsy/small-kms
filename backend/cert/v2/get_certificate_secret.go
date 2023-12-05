package cert

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetCertificateSecret implements admin.ServerInterface.
func (s *CertServer) GetCertificateSecret(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	cert, err := GetCertificateInternal(c, namespaceProvider, namespaceId, id)
	if err != nil {
		return err
	}
	if cert.GetJsonWebKey().Extractable == nil || !*cert.GetJsonWebKey().Extractable {
		return fmt.Errorf("%w: certificate key not extractable", base.ErrResponseStatusBadRequest)
	}

	req := new(certmodels.CertificateSecretRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	jwk, err := cloudkey.NewEphemeralECDHJwk(s.cryptoProvider)
	if err != nil {
		return err
	}

	jweBuilder := cloudkey.JWEAes256GcmEncBuilder{}
	if err := jweBuilder.SetEcdhEsKeyAgreement(jwk, &req.Jwk); err != nil {
		return err
	}

	sClient := kv.GetAzKeyVaultService(c).AzSecretsClient()
	sid := cert.KeyVaultSecretID()
	if sid == "" {
		return fmt.Errorf("%w: certificate secret not extractable", base.ErrResponseStatusBadRequest)
	}
	secretId := azsecrets.ID(sid)
	resp, err := sClient.GetSecret(c, secretId.Name(), secretId.Version(), nil)
	if err != nil {
		return err
	}
	pemBytes := []byte(*resp.Value)
	payload, err := jweBuilder.Seal(pemBytes)
	if err != nil {
		return nil
	}

	return c.JSON(200, &certmodels.CertificateSecretResult{
		Payload: payload,
	})
}
