package cert

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// PutCertificatePolicy implements ServerInterface.
func (*CertServer) PutCertificatePolicy(ec echo.Context, nsProvider models.NamespaceProvider, nsID string, ID string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	req := new(certmodels.CreateCertificatePolicyRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := ns.ValidateID(ID); err != nil {
		return err
	}

	doc := &CertPolicyDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: nsProvider,
				NamespaceID:       nsID,
				ResourceProvider:  models.ResourceProviderCertPolicy,
			},
			ID: ID,
		},
	}
	if err := doc.init(req); err != nil {
		return err
	}

	resp, err := resdoc.GetDocService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
