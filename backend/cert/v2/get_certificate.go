package cert

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetCertificate implements admin.ServerInterface.
func (*CertServer) GetCertificate(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string, params admin.GetCertificateParams) error {
	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	certDoc, err := getCertificateInternal(c, namespaceProvider, namespaceId, id)
	if err != nil {
		return err
	}

	includeJwk := false
	if params.IncludeJwk != nil {
		includeJwk = *params.IncludeJwk
	}
	model := certDoc.ToModel(includeJwk)
	return c.JSON(http.StatusOK, model)
}

func getCertificateInternal(c ctx.RequestContext, namespaceProvider models.NamespaceProvider, namespaceId string, id string) (*CertDoc, error) {
	certDoc := &CertDoc{}
	if err := resdoc.GetDocService(c).Read(c, resdoc.NewDocIdentifier(namespaceProvider, namespaceId, models.ResourceProviderCert, id), certDoc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return nil, base.ErrResponseStatusNotFound
		}
		return nil, err
	}
	return certDoc, nil
}
