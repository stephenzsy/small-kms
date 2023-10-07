package certtemplate

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateTemplateService interface {
	GetCertificateTemplate(common.ServiceContext, common.Identifier) (*models.CertificateTemplate, error)
	PutCertificateTemplate(common.ServiceContext, common.Identifier, models.CertificateTemplateParameters) (*models.CertificateTemplate, error)
	ListCertificateTemplates(common.ServiceContext) ([]*models.CertificateTemplateRef, error)
	WithCertificateTemplateContext(common.ServiceContext, common.Identifier) (common.ServiceContext, error)
}

type certTmplService struct {
}

type certTmplContextKey string

const (
	certTmplContextKeyDefault certTmplContextKey = "certTmplContext"
)

// WithCertificateTemplateContext implements CertificateTemplateService.
func (*certTmplService) WithCertificateTemplateContext(c common.ServiceContext, templateID common.Identifier) (common.ServiceContext, error) {
	ctc := newCertificateTemplateContext(templateID)
	return context.WithValue(c, certTmplContextKeyDefault, ctc), nil
}

func NewCertificateTemplateService() CertificateTemplateService {
	return &certTmplService{}
}

func GetCertificateTemplateContext(c common.ServiceContext) CertificateTemplateContext {
	if ctc, ok := c.Value(certTmplContextKeyDefault).(CertificateTemplateContext); ok {
		return ctc
	}
	return nil
}
