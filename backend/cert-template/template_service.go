package certtemplate

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
)

type certTmplContextKey string

const (
	certTmplContextKeyDefault certTmplContextKey = "certTmplContext"
)

// WithCertificateTemplateContext implements CertificateTemplateService.
func WithCertificateTemplateContext(c RequestContext, templateID common.Identifier) (RequestContext, error) {
	if !templateID.IsValid() {
		return nil, fmt.Errorf("%w:invalid template ID:%s", common.ErrStatusBadRequest, templateID)
	}
	ctc := newCertificateTemplateService(templateID)
	return common.RequestContextWithValue(c, certTmplContextKeyDefault, ctc), nil
}

func GetCertificateTemplateContext(c RequestContext) CertificateTemplateService {
	if ctc, ok := c.Value(certTmplContextKeyDefault).(CertificateTemplateService); ok {
		return ctc
	}
	return nil
}
