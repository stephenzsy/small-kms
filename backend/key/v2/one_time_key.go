package key

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type OneTimeKeyDoc struct {
	resdoc.ResourceDoc
	JWK      cloudkey.JsonWebKey `json:"jwk"`
	NotAfter resdoc.NumericDate  `json:"exp"`
	Created  resdoc.NumericDate  `json:"iat"`
}

func (doc *OneTimeKeyDoc) ToModel() *keymodels.OneTimeKey {
	if doc == nil {
		return nil
	}
	m := &keymodels.OneTimeKey{
		Jwk: *doc.JWK.PublicJWK(),
		Exp: doc.NotAfter,
		Iat: doc.Created,
	}
	return m
}

// CreateOneTimeKey implements admin.ServerInterface.
func (s *KeyAdminServer) CreateOneTimeKey(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string) error {
	c := ec.(ctx.RequestContext)

	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	keyId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	jwk, err := cloudkey.NewEphemeralECDHJwk(s.cryptoProvider)
	if err != nil {
		return err
	}
	now := time.Now().Truncate(time.Second)
	doc := &OneTimeKeyDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: namespaceProvider,
				NamespaceID:       namespaceId,
				ResourceProvider:  models.ResourceProviderOneTimeKey,
			},
			ID: keyId.String(),
		},
		JWK:      *jwk,
		Created:  *jwt.NewNumericDate(now),
		NotAfter: *jwt.NewNumericDate(now.Add(1 * time.Hour)),
	}
	doc.JWK.KeyID = doc.ID

	_, err = resdoc.GetDocService(c).Create(c, doc, nil)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, doc.ToModel())
}

func ReadOneTimeKey(c ctx.RequestContext, namespaceProvider models.NamespaceProvider, namespaceId string, keyID string) (*cloudkey.JsonWebKey, error) {
	c = c.Elevate()

	doc := &OneTimeKeyDoc{}

	docSvc := resdoc.GetDocService(c)
	err := docSvc.Read(c, resdoc.NewDocIdentifier(namespaceProvider, namespaceId, models.ResourceProviderOneTimeKey, keyID), doc, nil)
	if err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return nil, base.ErrResponseStatusNotFound
		}
		return nil, err
	}

	_, err = docSvc.Delete(c, doc.Identifier(), nil)
	if err != nil {
		return nil, err
	}

	if doc.NotAfter.Time.Before(time.Now()) {
		return nil, base.ErrResponseStatusNotFound
	}

	return &doc.JWK, err
}
