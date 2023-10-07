package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
)

type CertificateTemplateContext interface {
	GetCertificateTemplateDoc(common.ServiceContext) (*CertificateTemplateDoc, error)
}

type certTmplContext struct {
	templateID common.Identifier
}

// GetCertificateTemplateDoc implements CertificateTemplateContext.
func (ctc *certTmplContext) GetCertificateTemplateDoc(c common.ServiceContext) (*CertificateTemplateDoc, error) {
	return getCertificateTemplateDoc(c, ctc.templateID)
}

func newCertificateTemplateContext(templateID common.Identifier) CertificateTemplateContext {
	return &certTmplContext{
		templateID: templateID,
	}
}
