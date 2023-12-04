package cert

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GenerateCertificate implements ServerInterface.
func (*CertServer) GenerateCertificate(ec echo.Context,
	nsProvider models.NamespaceProvider, nsID string, policyID string) (err error) {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	var policy *CertPolicyDoc
	policy, err = GetCertificatePolicyInternal(c, nsProvider, nsID, policyID)
	if err != nil {
		return err
	}

	if !policy.AllowGenerate {
		return base.ErrResponseStatusBadRequest
	}

	var certDoc CertDocumentPending
	if policy.IssuerPolicy.NamespaceProvider == models.NamespaceProviderExternalCA {
		pending := &certDocACME{}
		pending.init(c, nsProvider, nsID, policy, nil)
		certDoc = pending
	} else {
		pending := &certDocInternal{}
		pending.init(c, nsProvider, nsID, policy, nil)
		certDoc = pending
	}

	docSvc := resdoc.GetDocService(c)
	if certAuthorized, err := certDoc.Authorize(c); err != nil {
		return err
	} else if !certAuthorized {
		if _, err := docSvc.Create(c, certDoc, nil); err != nil {
			return err
		}
		return c.JSON(http.StatusAccepted, certDoc.ToModel(true))
	}

	var csr CertCSR
	if nsProvider != models.NamespaceProviderRootCA {
		csr, err = certDoc.GetCertificateRequest(c, true)
		if err != nil {
			return err
		}
	}

	der, err := certDoc.CreateCertificate(c, csr)
	if err != nil {
		return err
	}
	if err := certDoc.CollectSignedCertificate(c, der); err != nil {
		return err
	}
	resp, err := docSvc.Create(c, certDoc, nil)
	if err != nil {
		return err
	}
	return c.JSON(resp.RawResponse.StatusCode, certDoc.ToModel(true))
}
