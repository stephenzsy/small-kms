package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type RequestContext = ctx.RequestContext

type CertificateTemplateService interface {
	GetCertificateTemplateLocator(RequestContext) shared.ResourceLocator
}

type certTmplContext struct {
	templateID common.Identifier
}

// GetCertificateTemplateLocator implements CertificateTemplateContext.
func (ctc *certTmplContext) GetCertificateTemplateLocator(c RequestContext) shared.ResourceLocator {
	nsID := ns.GetNamespaceContext(c).GetID()
	return getCertificateTemplateDocLocator(nsID, ctc.templateID)
}

func newCertificateTemplateService(templateID common.Identifier) CertificateTemplateService {
	return &certTmplContext{
		templateID: templateID,
	}
}
