package certtemplate

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
)

type certTmplContextKey string

const (
	certTmplContextKeyDefault certTmplContextKey = "certTmplContext"
)

// WithCertificateTemplateContext implements CertificateTemplateService.
func WithCertificateTemplateContext(c common.ServiceContext, templateID common.Identifier) (common.ServiceContext, error) {
	if !templateID.IsValid() {
		return nil, fmt.Errorf("%w:invalid template ID:%s", common.ErrStatusBadRequest, templateID)
	}
	ctc := newCertificateTemplateContext(templateID)
	return context.WithValue(c, certTmplContextKeyDefault, ctc), nil
}

func GetCertificateTemplateContext(c common.ServiceContext) CertificateTemplateContext {
	if ctc, ok := c.Value(certTmplContextKeyDefault).(CertificateTemplateContext); ok {
		return ctc
	}
	return nil
}
