package cert

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// PutCertificatePolicyIssuer implements admin.ServerInterface.
func (*CertServer) PutCertificatePolicyIssuer(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	if namespaceProvider != models.NamespaceProviderRootCA && namespaceProvider != models.NamespaceProviderIntermediateCA {
		return base.ErrResponseStatusBadRequest
	}

	req := new(models.LinkRefFields)
	if err := c.Bind(req); err != nil {
		return err
	}

	linkTo, err := resdoc.ParseIdentifier(req.LinkTo)
	if err != nil {
		return fmt.Errorf("%w: %w", base.ErrResponseStatusBadRequest, err)
	}
	if linkTo.NamespaceProvider != namespaceProvider || linkTo.NamespaceID != namespaceId || linkTo.ResourceProvider != models.ResourceProviderCert {
		return fmt.Errorf("%w: %s", base.ErrResponseStatusBadRequest, "invalid linkTo")
	}
	certDoc, err := GetCertificateInternal(c, namespaceProvider, namespaceId, linkTo.ID)
	if err != nil {
		return err
	}
	if certDoc.GetStatus() != certmodels.CertificateStatusIssued {
		return fmt.Errorf("%w: %s", base.ErrResponseStatusBadRequest, "certificate is not issued")
	}
	doc := &resdoc.LinkResourceDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: namespaceProvider,
				NamespaceID:       namespaceId,
				ResourceProvider:  models.ResourceProviderLink,
			},
			ID: getPolicyIssuerCertLinkID(id),
		},
		LinkTo:       linkTo,
		LinkProvider: models.LinkProviderCAPolicyIssuerCertificate,
	}
	resp, err := resdoc.GetDocService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
