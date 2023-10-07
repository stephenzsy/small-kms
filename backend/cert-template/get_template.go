package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/profile"
)

// PutCertificateTemplate implements CertificateTemplateService.
func (s *certTmplService) GetCertificateTemplate(c common.ServiceContext,
	templateID models.Identifier) (*models.CertificateTemplate, error) {

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	pcs := profile.GetProfileContextService(c)
	nsID := pcs.GetResourceDocNsID()

	doc := CertificateTemplateDoc{}
	if err := kmsdoc.Read(c, nsID, kmsdoc.NewDocIdentifier(kmsdoc.DocKindCertificateTemplate, templateID), &doc); err != nil {
		return nil, err
	}

	return doc.toModel(), nil
}
