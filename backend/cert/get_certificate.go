package cert

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func NewCertificateID(certId common.Identifier) models.ResourceID {
	return common.NewIdentifierWithKind(models.ResourceKindCert, certId)
}

func NewLatestCertificateForTemplateID(certId common.Identifier) models.ResourceID {
	return common.NewIdentifierWithKind(models.ResourceKindLatestCertForTemplate, certId)
}

type NamespaceID = models.NamespaceID
type ResourceID = models.ResourceID
type Identifier = common.Identifier

func getCrossNsReferencedTemplateIdentifier(referencedNamespaceID NamespaceID, templateIdentifier Identifier) Identifier {
	uuidValue := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("%s/%s", referencedNamespaceID.String(), templateIdentifier.String())))
	return common.UUIDIdentifier(uuidValue)
}

func GetCertificate(c RequestContext, certificateId common.Identifier, params models.GetCertificateParams) (*models.CertificateInfoComposed, error) {
	var certDocLocator models.ResourceLocator
	if certificateId.IsUUID() {
		nsID := ns.GetNamespaceContext(c).GetID()
		certDocLocator = models.NewResourceLocator(nsID, NewCertificateID(certificateId))
	} else if certificateId.String() == "latest" {
		if params.TemplateId.IsNilOrEmpty() || !params.TemplateId.IsValid() {
			return nil, fmt.Errorf("%w: invalid or empty template ID: %s", common.ErrStatusBadRequest, params.TemplateId)
		}
		if !params.TemplateNamespaceId.IsNilOrEmpty() {
			if !params.TemplateNamespaceId.IsValid() || params.TemplateNamespaceKind == nil || *params.TemplateNamespaceKind != models.NamespaceKindGroup {
				return nil, fmt.Errorf("%w: invalid template namespace ID: %s", common.ErrStatusBadRequest, params.TemplateNamespaceId)
			}
			nsID := ns.GetNamespaceContext(c).GetID()

			certDocLocator = models.NewResourceLocator(nsID,
				NewLatestCertificateForTemplateID(getCrossNsReferencedTemplateIdentifier(
					models.NewNamespaceID(*params.TemplateNamespaceKind, *params.TemplateId),
					*params.TemplateId)))
		} else {
			// same namespace
			nsID := ns.GetNamespaceContext(c).GetID()
			certDocLocator = models.NewResourceLocator(nsID, NewLatestCertificateForTemplateID(*params.TemplateId))
		}
	} else {
		return nil, fmt.Errorf("%w: invalid certificate ID: %s", common.ErrStatusBadRequest, certificateId)
	}

	certDoc := &CertDoc{}
	err := kmsdoc.Read(c, certDocLocator, certDoc)
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
