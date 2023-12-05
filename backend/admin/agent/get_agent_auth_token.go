package agentadmin

import (
	"github.com/labstack/echo/v4"
	agentauth "github.com/stephenzsy/small-kms/backend/agent/auth"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	cloudkeyx "github.com/stephenzsy/small-kms/backend/cloud/key/x"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

// GetAgentAuthToken implements admin.ServerInterface.
func (*AgentAdminServer) GetAgentAuthToken(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}
	instance, err := getAgentInstanceInternal(c, namespaceId, id)
	if err != nil {
		return err
	}
	if instance.JwtVerfyKeyID == "" {
		return base.ErrResponseStatusBadRequest
	}

	identity := auth.GetAuthIdentity(c)
	keyDoc, err := key.GetKeyInternal(c, models.NamespaceProviderServicePrincipal, namespaceId, instance.JwtVerfyKeyID)
	if err != nil {
		return err
	}

	ck := cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, kv.GetAzKeyVaultService(c).AzKeysClient(), keyDoc.JsonWebKey.KeyID, cloudkey.SignatureAlgorithmES384, false, keyDoc.PublicKey())

	accessToken, _, err := agentauth.NewSignedAgentAuthJWT(cloudkeyx.NewJWTSigningMethod(cloudkey.SignatureAlgorithmES384), identity.ClientPrincipalID().String(), instance.Endpoint, ck)
	if err != nil {
		return err
	}

	return c.JSON(200, &agentmodels.AgentAuthResult{
		AccessToken: accessToken,
	})
}
