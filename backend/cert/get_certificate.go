package cert

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func NewCertificateID(certId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindCert, certId)
}

func NewLatestCertificateForTemplateID(templateId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate, templateId)
}

func getCrossNsReferencedTemplateIdentifier(referencedNamespaceID shared.NamespaceIdentifier, templateIdentifier shared.Identifier) shared.Identifier {
	uuidValue := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("%s/%s", referencedNamespaceID.String(), templateIdentifier.String())))
	return shared.UUIDIdentifier(uuidValue)
}

func ReadCertDocByLocator(c context.Context, locator shared.ResourceLocator) (*CertDoc, error) {
	certDoc := &CertDoc{}
	err := kmsdoc.Read(c, locator, certDoc)
	return certDoc, err
}

func ApiGetCertificate(c RequestContext, certificateId shared.Identifier, params models.GetCertificateParams) error {
	cert, err := GetCertificate(c, certificateId, params)
	if err != nil {
		return err
	}
	return c.JSON(200, cert)
}

func GetCertificate(c RequestContext, certificateId shared.Identifier, params models.GetCertificateParams) (*models.CertificateInfoComposed, error) {
	var certDocLocator shared.ResourceLocator
	if certificateId.IsUUID() {
		nsID := ns.GetNamespaceContext(c).GetID()
		certDocLocator = shared.NewResourceLocator(nsID, NewCertificateID(certificateId))
	} else {
		return nil, fmt.Errorf("%w: invalid certificate ID: %s", common.ErrStatusBadRequest, certificateId)
	}

	certDoc, err := ReadCertDocByLocator(c, certDocLocator)
	if err != nil {
		return nil, err
	}
	m := certDoc.toModel()

	if params.IncludeCertificate != nil && (*params.IncludeCertificate == models.IncludePEM || *params.IncludeCertificate == models.IncludeJWK) {
		// fetch cert from blob
		pemBlob, err := certDoc.fetchCertificatePEMBlob(c)
		if err != nil {
			return m, err
		}
		m.Pem = utils.ToPtr(string(pemBlob))
		switch *params.IncludeCertificate {
		case models.IncludePEM:
			m.Pem = utils.ToPtr(string(pemBlob))
		case models.IncludeJWK:
		}
		// attach blob
	}

	return m, nil
}
