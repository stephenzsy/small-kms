package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

// PutCertificateTemplate implements CertificateTemplateService.
func (s *certTmplService) PutCertificateTemplate(c common.ServiceContext,
	templateID models.Identifier,
	req models.CertificateTemplateParameters) (*models.CertificateTemplate, error) {

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	doc, err := validatePutRequest(c, templateID, req)
	if err != nil {
		return nil, err
	}

	err = kmsdoc.Upsert(c, doc)
	if err != nil {
		return nil, err
	}

	return doc.toModel(), nil
}
