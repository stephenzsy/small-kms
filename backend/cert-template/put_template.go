package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

// PutCertificateTemplate implements CertificateTemplateService.
func (s *certTmplService) PutCertificateTemplate(c common.ServiceContext,
	templateID models.Identifier,
	req models.CertificateTemplateParameters) (*models.CertificateTemplateParameters, error) {
	panic("unimplemented")
}
