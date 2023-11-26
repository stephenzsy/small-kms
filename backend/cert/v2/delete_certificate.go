package cert

import (
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// DeleteCertificate implements admin.ServerInterface.
func (*CertServer) DeleteCertificate(ec echo.Context,
	namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	doc, err := GetCertificateInternal(c, namespaceProvider, namespaceId, id)
	if err != nil {
		if errors.Is(err, base.ErrResponseStatusNotFound) {
			return c.NoContent(http.StatusNoContent)
		}
		return err
	}

	c = c.Elevate()
	if err = doc.cleanupKeyVault(c); err != nil {
		return err
	}
	if doc.Status == certmodels.CertificateStatusPending || doc.NotAfter.Time.Before(time.Now()) {

		// delete document
		_, err := resdoc.GetDocService(c).Delete(c, doc.Identifier(), &azcosmos.ItemOptions{
			IfMatchEtag: doc.ETag,
		})
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	}
	// will be put to deactivated state
	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendSet("/status", certmodels.CertificateStatusDeactivated)
	patchOps.AppendSet("/deleted", time.Now().UTC())
	resp, err := resdoc.GetDocService(c).Patch(c, doc, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: doc.ETag,
	})
	if err != nil {
		return err
	}
	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel(true))
}
