package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateTemplateService interface {
	PutCertificateTemplate(common.ServiceContext, models.Identifier, models.CertificateTemplateParameters) (*models.CertificateTemplate, error)
	ListCertificateTemplates(common.ServiceContext) ([]*models.CertificateTemplateRef, error)
}

type certTmplService struct {
}

func NewCertificateTemplateService() CertificateTemplateService {
	return &certTmplService{}
}
