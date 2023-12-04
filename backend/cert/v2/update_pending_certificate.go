package cert

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"golang.org/x/crypto/acme"
)

// UpdatePendingCertificate implements admin.ServerInterface.
func (*CertServer) UpdatePendingCertificate(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	req := new(certmodels.UpdatePendingCertificateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	certDoc := new(certDocACME)
	if err := readCertDocInternal(c, namespaceProvider, namespaceId, id, certDoc); err != nil {
		return err
	}

	if certDoc.Status != certmodels.CertificateStatusPendingAuthorization {
		return base.ErrResponseStatusNotFound
	}

	if err := certDoc.restore(c); err != nil {
		return err
	}

	switch {
	case req.AcmeAcceptChallenge != "":
		_, err := certDoc.acmeClient.Accept(c, &acme.Challenge{
			URI: req.AcmeAcceptChallenge,
		})
		if err != nil {
			return err
		}
	case req.AcmeOrderCertificate != nil && *req.AcmeOrderCertificate:

		csr, err := certDoc.GetCertificateRequest(c, false)
		if err != nil {
			return err
		}

		der, err := certDoc.CreateCertificate(c, csr)
		if err != nil {
			return err
		}

		if err := certDoc.CollectSignedCertificate(c, der); err != nil {
			return err
		}
		docSvc := resdoc.GetDocService(c)
		if _, err := docSvc.Upsert(c, certDoc, &azcosmos.ItemOptions{
			IfMatchEtag: certDoc.ETag,
		}); err != nil {
			return err
		}
	}

	model := certDoc.ToModel(true)
	return c.JSON(http.StatusOK, model)
}
