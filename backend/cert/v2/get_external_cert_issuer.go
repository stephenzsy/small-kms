package cert

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

func (*CertServer) GetExternalCertificateIssuer(ec echo.Context, namespaceId string, issuerID string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	doc, err := getExternalCertificateIssuerInternal(c, namespaceId, issuerID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, doc.ToModel())
}

func getExternalCertificateIssuerInternal(c ctx.RequestContext, namespaceId string, issuerID string) (*CertIssuerDoc, error) {
	doc := &CertIssuerDoc{}
	docSvc := resdoc.GetDocService(c)
	err := docSvc.Read(c,
		resdoc.NewDocIdentifier(models.NamespaceProviderExternalCA, namespaceId, models.ResourceProviderCertExternalIssuer, issuerID), doc, nil)
	if err != nil && errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
		return nil, base.ErrResponseStatusNotFound
	}
	return doc, err
}
