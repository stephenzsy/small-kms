package cert

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
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

	if params.Pending != nil && *params.Pending {
		return getCertificatePending(c, namespaceProvider, namespaceId, id)
	}

	certDoc, err := GetCertificateInternal(c, namespaceProvider, namespaceId, id)
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

func GetCertificateInternal(c context.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) (CertDocument, error) {
	certDoc := &certDocBase{}
	err := readCertDocInternal(c, namespaceProvider, namespaceId, id, certDoc)
	if err != nil {
		return nil, err
	}
	return certDoc, nil
}

func readCertDocInternal[T resdoc.ResourceDocument](c context.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string, certDoc T) error {
	if err := resdoc.GetDocService(c).Read(c, resdoc.NewDocIdentifier(namespaceProvider, namespaceId, models.ResourceProviderCert, id), certDoc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: certificate ID: %s", base.ErrResponseStatusNotFound, id)
		}
		return err
	}
	return nil
}

func getCertificatePending(c ctx.RequestContext, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	certDoc := &certDocACME{}
	if err := readCertDocInternal(c, namespaceProvider, namespaceId, id, certDoc); err != nil {
		return err
	}
	if err := certDoc.restore(c); err != nil {
		return err
	}
	order, err := certDoc.acmeClient.GetOrder(c, certDoc.OrderURL)
	if err != nil {
		return err
	}
	authorizations := make([]certmodels.CertificatePendingAcmeAuthorization, len(order.AuthzURLs))

	for i, url := range order.AuthzURLs {
		a, err := certDoc.acmeClient.GetAuthorization(c, url)
		if err != nil {
			return err
		}
		authorizations[i] = certmodels.CertificatePendingAcmeAuthorization{
			Challenges: make([]certmodels.CertificatePendingAcmeChallenge, len(a.Challenges)),
			Status:     a.Status,
			URL:        a.URI,
		}
		for j, ch := range a.Challenges {
			if ch.Type == "dns-01" {
				record, err := certDoc.acmeClient.DNS01ChallengeRecord(ch.Token)
				if err != nil {
					return err
				}
				authorizations[i].Challenges[j] = certmodels.CertificatePendingAcmeChallenge{
					Type:      ch.Type,
					DNSRecord: record,
					URL:       ch.URI,
				}
			} else {
				authorizations[i].Challenges[j] = certmodels.CertificatePendingAcmeChallenge{
					Type: ch.Type,
					URL:  ch.URI,
				}
			}
		}
	}
	model := certDoc.ToModel(true)
	return c.JSON(http.StatusOK, model)
}
