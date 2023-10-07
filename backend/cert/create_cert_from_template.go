package cert

import (
	"fmt"

	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

var (
	ErrInvalidContext = fmt.Errorf("invalid context")
)

type certificateContext struct {
}

func createCertContextFromTemplate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*certificateContext, error) {

	ctc := ct.GetCertificateTemplateContext(c)
	_, err := ctc.GetCertificateTemplateDoc(c)
	if err != nil {
		return nil, err
	}

	panic("unimplemented")
}

func issueCertificate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*CertDoc, error) {

	_ = ct.GetCertificateTemplateContext(c)
	_, ok := c.Value(certContext).(*certificateContext)
	if !ok {
		return nil, fmt.Errorf("%w: no cert context to issue", ErrInvalidContext)
	}

	panic("unimplemented")
}
