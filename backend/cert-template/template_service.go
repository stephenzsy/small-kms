package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateTemplateService interface {
	GetCertificateTemplate(common.ServiceContext, common.Identifier) (*models.CertificateTemplate, error)
	PutCertificateTemplate(common.ServiceContext, common.Identifier, models.CertificateTemplateParameters) (*models.CertificateTemplate, error)
	ListCertificateTemplates(common.ServiceContext) ([]*models.CertificateTemplateRef, error)
}

type certTmplService struct {
}

func NewCertificateTemplateService() CertificateTemplateService {
	return &certTmplService{}
}
