package key

import (
	"crypto/sha512"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GenerateKey implements admin.ServerInterface.
func (*KeyAdminServer) GenerateKey(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	policy, err := getKeyPolicyInternal(c, namespaceProvider, namespaceId, id)
	if err != nil {
		return err
	}

	doc := &keyGenerateDoc{
		KeyDoc: KeyDoc{
			ResourceDoc: resdoc.ResourceDoc{
				PartitionKey: resdoc.PartitionKey{
					NamespaceProvider: namespaceProvider,
					NamespaceID:       namespaceId,
					ResourceProvider:  models.ResourceProviderKey,
				},
			},
		},
	}
	err = doc.init(namespaceProvider, namespaceId, policy)
	if err != nil {
		return err
	}
	params, err := doc.getAzCreateKeyParams()
	if err != nil {
		return err
	}
	c = c.Elevate()
	result, err := kv.GetAzKeyVaultService(c).AzKeysClient().CreateKey(c, doc.keyVaultStoreName, params, nil)
	if err != nil {
		return err
	}
	doc.KeyID = string(*result.Key.KID)
	doc.Created.Time = *result.Attributes.Created
	if result.Attributes.NotBefore != nil {
		doc.NotBefore = jwt.NewNumericDate(*result.Attributes.NotBefore)
	}
	if result.Attributes.Expires != nil {
		doc.NotAfter = jwt.NewNumericDate(*result.Attributes.Expires)
	}
	doc.N = result.Key.N
	doc.E = result.Key.E
	doc.X = result.Key.X
	doc.Y = result.Key.Y
	doc.Exportable = *result.Attributes.Exportable
	doc.Status = keymodels.KeyStatusActive

	digest := sha512.New384()
	doc.JsonWebKey.Digest(digest)

	if doc.NotBefore != nil {
		if m, _ := doc.NotBefore.MarshalJSON(); m != nil {
			digest.Write(m)
		}
	}
	if doc.NotAfter != nil {
		if m, _ := doc.NotAfter.MarshalJSON(); m != nil {
			digest.Write(m)
		}
	}
	doc.Checksum = digest.Sum(nil)
	resp, err := resdoc.GetDocService(c).Create(c, doc, nil)
	if err != nil {
		return err
	}
	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel(true))
}
