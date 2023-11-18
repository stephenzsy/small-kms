package cert

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"

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
	policy, err = getCertificatePolicyInternal(c, nsProvider, nsID, policyID)
	if err != nil {
		return err
	}

	certDoc := &certDocSelfSignedGeneratePending{
		CertDoc: CertDoc{
			ResourceDoc: resdoc.ResourceDoc{
				PartitionKey: resdoc.PartitionKey{
					NamespaceProvider: nsProvider,
					NamespaceID:       nsID,
					ResourceProvider:  models.ResourceProviderCert,
				},
			},
		},
	}

	c = c.Elevate()
	err = certDoc.init(c, nsProvider, nsID, policy)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cancelErr := certDoc.cancel(c); cancelErr != nil {
				err = fmt.Errorf("%w+%w", err, cancelErr)
			}
		}
		resp, docCreateErr := resdoc.GetDocService(c).Create(c, certDoc, nil)
		if docCreateErr != nil {
			if err == nil {
				err = docCreateErr
			} else {
				err = fmt.Errorf("%w+%w", err, docCreateErr)
			}
		}
		if err != nil {
			return
		}
		err = c.JSON(resp.RawResponse.StatusCode, certDoc.ToModel())
	}()
	var signed []byte
	signed, err = x509.CreateCertificate(rand.Reader,
		certDoc.templateX509Cert,
		certDoc.issuerX509Cert,
		certDoc.publicKey,
		certDoc.signer)
	if err != nil {
		return
	}
	err = certDoc.collectSignedCert(signed)
	if err != nil {
		return
	}
	return
}
