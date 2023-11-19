package cert

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetCertificatePolicyIssuer implements admin.ServerInterface.
func (*CertServer) GetCertificatePolicyIssuer(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, policyID string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	if namespaceProvider != models.NamespaceProviderRootCA && namespaceProvider != models.NamespaceProviderIntermediateCA {
		return base.ErrResponseStatusBadRequest
	}

	doc, err := getPolicyIssuerCertInternal(c, namespaceProvider, namespaceId, policyID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, doc.ToModel())
}

func getPolicyIssuerCertInternal(c ctx.RequestContext, namespaceProvider models.NamespaceProvider, namespaceId string, policyID string) (*resdoc.LinkResourceDoc, error) {
	doc := new(resdoc.LinkResourceDoc)
	if err := resdoc.GetDocService(c).Read(c, resdoc.NewDocIdentifier(
		namespaceProvider, namespaceId, models.ResourceProviderLink, getPolicyIssuerCertLinkID(policyID)), doc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w, policyID: %s", base.ErrResponseStatusNotFound, policyID)
		}
		return nil, err
	}
	return doc, nil
}

func getPolicyIssuerCertLinkID(policyID string) string {
	return fmt.Sprintf("issuer-cert-%s", policyID)
}
