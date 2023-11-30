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

func GetCertificateInternal(c context.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) (*CertDoc, error) {
	certDoc := &CertDoc{}
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
	certDoc := &certDocExternalACMEPending{}
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
	challenges := make([]certmodels.CertificatePendingAcmeChallenge, 0)
	for _, url := range order.AuthzURLs {
		a, err := certDoc.acmeClient.GetAuthorization(c, url)
		if err != nil {
			return err
		}
		for _, ch := range a.Challenges {
			if ch.Type == "dns-01" {
				record, err := certDoc.acmeClient.DNS01ChallengeRecord(ch.Token)
				if err != nil {
					return err
				}
				challenges = append(challenges, certmodels.CertificatePendingAcmeChallenge{
					DNSRecord: record,
				})
			}
		}
	}
	model := certDoc.toModel(&certmodels.CertificatePendingAcme{
		Challenges: challenges,
	})
	return c.JSON(http.StatusOK, model)
}
