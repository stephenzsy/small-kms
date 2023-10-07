package cert

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/profile"
)

// CreateCertificateFromTemplate implements CertificateService.
func (*certService) CreateCertificateFromTemplate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*models.CertificateInfo, error) {

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	_ = profile.GetProfileContext(c)
	ctc := ct.GetCertificateTemplateContext(c)
	ctc.GetCertificateTemplateDoc(c)

	panic("unimplemented")
}
