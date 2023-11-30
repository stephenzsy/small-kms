package cert

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
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

	if policy.IssuerPolicy.NamespaceProvider == models.NamespaceProviderExternalCA {
		return generateExternal(c, nsProvider, nsID, policy)
	}

	certDoc := &certDocGeneratePending{
		certDocPending: certDocPending{
			CertDoc: CertDoc{
				ResourceDoc: resdoc.ResourceDoc{
					PartitionKey: resdoc.PartitionKey{
						NamespaceProvider: nsProvider,
						NamespaceID:       nsID,
						ResourceProvider:  models.ResourceProviderCert,
					},
				},
			},
		},
	}

	c = c.Elevate()
	err = certDoc.init(c, nsProvider, nsID, policy)
	defer func() {
		if err != nil {
			if cancelErr := certDoc.cleanupKeyVault(c); cancelErr != nil {
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
		err = c.JSON(resp.RawResponse.StatusCode, certDoc.ToModel(true))
	}()
	if err != nil {
		return err
	}

	var signed []byte
	signed, err = x509.CreateCertificate(rand.Reader,
		certDoc.templateX509Cert,
		certDoc.issuerX509Cert,
		certDoc.publicKey,
		certDoc.signer)
	if err != nil {
		return
	}
	err = certDoc.collectSignedCert(c, signed)
	if err != nil {
		return
	}
	return
}

func generateExternal(c ctx.RequestContext, nsProvider models.NamespaceProvider, nsID string, policyDoc *CertPolicyDoc) error {

	// load certificate issuer

	certDoc := &certDocExternalACMEPending{
		certDocPending: certDocPending{
			CertDoc: CertDoc{
				ResourceDoc: resdoc.ResourceDoc{
					PartitionKey: resdoc.PartitionKey{
						NamespaceProvider: nsProvider,
						NamespaceID:       nsID,
						ResourceProvider:  models.ResourceProviderCert,
					},
				},
			},
		},
	}

	err := certDoc.init(c, nsProvider, nsID, policyDoc)
	if err != nil {
		return err
	}

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Create(c, certDoc, nil)
	if err != nil {
		return err
	}

	order, err := certDoc.authorizeOrder(c)
	if err != nil {
		return err
	}
	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendSet("/orderUrl", order.URI)
	patchOps.AppendSet("/acmeStep", AcmeStepOrderCreated)
	patchOps.AppendSet("/status", certmodels.CertificateStatusPendingExternal)
	_, err = docSvc.Patch(c, certDoc, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: certDoc.ETag,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, certDoc.ToModel(true))
}
