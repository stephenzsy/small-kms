package cert

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateService interface {
	CreateCertificateFromTemplate(common.ServiceContext, models.IssueCertificateFromTemplateParams) (*models.CertificateInfo, error)
}

type certService struct {
}

func NewCertificateService() CertificateService {
	return &certService{}
}
