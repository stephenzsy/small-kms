package cert

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetCertificatePolicy implements ServerInterface.
func (*CertServer) GetCertificatePolicy(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {

	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	doc, err := getCertificatePolicyInternal(c, namespaceProvider, namespaceId, id)
	if err != nil {
		return err
	}
	return c.JSON(200, doc.ToModel())
}

func getCertificatePolicyInternal(c ctx.RequestContext, namespaceProvider models.NamespaceProvider, namespaceId string, id string) (*CertPolicyDoc, error) {
	doc := &CertPolicyDoc{}
	if err := resdoc.GetDocService(c).Read(c, resdoc.DocIdentifier{
		PartitionKey: resdoc.PartitionKey{
			NamespaceProvider: namespaceProvider,
			NamespaceID:       namespaceId,
			ResourceProvider:  models.ResourceProviderCertPolicy,
		},
		ID: id,
	}, doc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: certificate policy not found: %s", base.ErrResponseStatusNotFound, id)
		}
		return nil, err
	}
	return doc, nil
}
