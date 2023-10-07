package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type CertificateTemplateContext interface {
	GetCertificateTemplateLocator(common.ServiceContext) models.ResourceLocator
	GetCertificateTemplateDoc(common.ServiceContext) (*CertificateTemplateDoc, error)
}

type certTmplContext struct {
	templateID common.Identifier
}

// GetCertificateTemplateLocator implements CertificateTemplateContext.
func (ctc *certTmplContext) GetCertificateTemplateLocator(c common.ServiceContext) models.ResourceLocator {
	nsID := ns.GetNamespaceContext(c).GetID()
	return getCertificateTemplateDocLocator(nsID, ctc.templateID)
}

// GetCertificateTemplateDoc implements CertificateTemplateContext.
func (ctc *certTmplContext) GetCertificateTemplateDoc(c common.ServiceContext) (*CertificateTemplateDoc, error) {
	return GetCertificateTemplateDoc(c, ctc.GetCertificateTemplateLocator(c))
}

func newCertificateTemplateContext(templateID common.Identifier) CertificateTemplateContext {
	return &certTmplContext{
		templateID: templateID,
	}
}
