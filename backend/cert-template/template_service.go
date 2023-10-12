package certtemplate

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type certTmplContextKey string

const (
	certTmplContextKeyDefault certTmplContextKey = "certTmplContext"
)

// WithCertificateTemplateContext implements CertificateTemplateService.
func WithCertificateTemplateContext(c RequestContext, templateID shared.Identifier) (RequestContext, error) {
	if !templateID.IsValid() {
		return c, fmt.Errorf("%w:invalid template ID:%s", common.ErrStatusBadRequest, templateID)
	}
	ctc := newCertificateTemplateService(templateID)
	return c.WithSharedValue(certTmplContextKeyDefault, ctc), nil
}

func GetCertificateTemplateContext(c RequestContext) CertificateTemplateService {
	if ctc, ok := c.Value(certTmplContextKeyDefault).(CertificateTemplateService); ok {
		return ctc
	}
	return nil
}
