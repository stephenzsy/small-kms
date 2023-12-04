package cert

import (
	"context"
	"net"

	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/utils"
	"golang.org/x/crypto/acme"
)

type AcmeStep string

const (
	AcmeStepNone         AcmeStep = ""
	AcmeStepOrderCreated AcmeStep = "orderCreated"
)

type certDocACME struct {
	certDocPending
	ACMEStep AcmeStep `json:"acmeStep"`
	OrderURL string   `json:"orderUrl"`

	acmeClient *acme.Client
}

// CreateCertificate implements CertDocumentPending.
func (doc *certDocACME) CreateCertificate(c ctx.RequestContext, csr CertCSR) ([][]byte, error) {
	order, err := doc.acmeClient.GetOrder(c, doc.OrderURL)
	if err != nil {
		return nil, err
	}
	certs, _, err := doc.acmeClient.CreateOrderCert(c, order.FinalizeURL, csr.X509CSRBytes(), true)
	if err != nil {
		return nil, err
	}
	return certs, err
}

func (doc *certDocACME) init(c ctx.RequestContext,
	nsProvider models.NamespaceProvider, nsID string,
	pDoc *CertPolicyDoc, publicKey *cloudkey.JsonWebKey) (err error) {
	err = doc.certDocPending.init(c, nsProvider, nsID, pDoc, publicKey)
	if issuerDoc, err := pDoc.getExternalIssuer(c); err != nil {
		return err
	} else {
		doc.Issuer = issuerDoc.Identifier()
		doc.acmeClient, err = issuerDoc.ACMEClient(c)
		if err != nil {
			return err
		}
	}
	return err
}

func (doc *certDocACME) restore(c ctx.RequestContext) error {
	issuerDoc, error := getExternalCertificateIssuerInternal(c, doc.Issuer.NamespaceID,
		doc.Issuer.ID)
	if error != nil {
		return error
	}
	doc.acmeClient, error = issuerDoc.ACMEClient(c)
	return error
}

func (doc *certDocACME) Authorize(c context.Context) (bool, error) {
	authIds := make([]acme.AuthzID, 0, len(doc.SANs.DNSNames)+len(doc.SANs.IPAddresses))
	authIds = append(authIds, acme.DomainIDs(doc.SANs.DNSNames...)...)
	authIds = append(authIds, acme.IPIDs(utils.MapSlice(doc.SANs.IPAddresses, func(ip net.IP) string {
		return ip.String()
	})...)...)

	order, err := doc.acmeClient.AuthorizeOrder(c, authIds)
	if err != nil {
		return false, err
	}
	doc.ACMEStep = AcmeStepOrderCreated
	doc.OrderURL = order.URI
	doc.Status = certmodels.CertificateStatusPendingAuthorization
	return order.Status == acme.StatusValid, nil
}

var _ CertDocumentPending = (*certDocACME)(nil)
