package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type RequestContext = common.RequestContext

type CertificateTemplateService interface {
	GetCertificateTemplateLocator(RequestContext) models.ResourceLocator
	GetCertificateTemplateDoc(RequestContext) (*CertificateTemplateDoc, error)
}

type certTmplContext struct {
	templateID common.Identifier
}

// GetCertificateTemplateLocator implements CertificateTemplateContext.
func (ctc *certTmplContext) GetCertificateTemplateLocator(c RequestContext) models.ResourceLocator {
	nsID := ns.GetNamespaceContext(c).GetID()
	return getCertificateTemplateDocLocator(nsID, ctc.templateID)
}

// GetCertificateTemplateDoc implements CertificateTemplateContext.
func (ctc *certTmplContext) GetCertificateTemplateDoc(c RequestContext) (*CertificateTemplateDoc, error) {
	return GetCertificateTemplateDoc(c, ctc.GetCertificateTemplateLocator(c))
}

func newCertificateTemplateService(templateID common.Identifier) CertificateTemplateService {
	return &certTmplContext{
		templateID: templateID,
	}
}
