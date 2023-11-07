package secret

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GenerateSecret implements ServerInterface.
func (s *server) GenerateSecret(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, policyID base.ID) error {
	c := ec.(ctx.RequestContext)
	logger := log.Ctx(c)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	if err := ns.VerifyKeyVaultIdentifier(policyID); err != nil {
		return fmt.Errorf("%w: invalid secret policy identifier: %s", base.ErrResponseStatusBadRequest, err)
	}

	c = ns.WithNSContext(c, namespaceKind, namespaceId)
	pDoc, err := readSecretPolicyDoc(c, policyID)
	if err != nil {
		return wrapSecretPolicyNotFoundError(err, policyID)
	}
	if pDoc.Mode != SecretGenerateModeServerGeneratedRandom {
		return fmt.Errorf("%w: invalid secret policy mode to generate secret: %s", base.ErrResponseStatusBadRequest, pDoc.Mode)
	}

	doc := &SecretDoc{}
	doc.init(namespaceKind, namespaceId, pDoc)

	// create secrets
	secretBytes := make([]byte, *pDoc.RandomLength)
	if _, err := cryptorand.Read(secretBytes); err != nil {
		return err
	}
	var encodedSecret string
	switch *pDoc.RandomCharacterClass {
	case SecretRandomCharClassBase64RawURL:
		encodedSecret = base64.RawURLEncoding.EncodeToString(secretBytes)
	default:
		encodedSecret = string(secretBytes)
	}

	resp, err := kv.GetAzKeyVaultService(c).AzSecretsClient().SetSecret(c, doc.KeyVaultStore.Name, azsecrets.SetSecretParameters{
		Value:       &encodedSecret,
		ContentType: to.Ptr("text/plain"),
		SecretAttributes: &azsecrets.SecretAttributes{
			Enabled: to.Ptr(true),
		},
	}, nil)
	if err != nil {
		return err
	}
	logger.Info().Str("secretId", string(*resp.ID)).Msgf("secret created in keyvault")

	doc.Version = resp.ID.Version()
	doc.KeyVaultStore.ID = string(*resp.ID)
	if resp.Attributes.Created != nil {
		doc.Created = *jwt.NewNumericDate(*resp.Attributes.Created)
	}
	if resp.Attributes.Expires != nil {
		doc.NotAfter = jwt.NewNumericDate(*resp.Attributes.Expires)
	}
	if resp.Attributes.NotBefore != nil {
		doc.NotBefore = jwt.NewNumericDate(*resp.Attributes.NotBefore)
	}
	docSvc := base.GetAzCosmosCRUDService(c)
	err = docSvc.Create(c, doc, nil)
	if err != nil {
		return err
	}

	m := &Secret{}
	doc.PopulateModel(m)
	return c.JSON(200, m)
}
