package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type certServiceContextKey string

const (
	certContext certServiceContextKey = "certContext"
)

// IssueCertificateFromTemplate implements CertificateService.
func IssueCertificateFromTemplate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*models.CertificateInfoComposed, error) {

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	certCtx, err := createCertContextFromTemplate(c, params)
	if err != nil {
		return nil, err
	}
	c = context.WithValue(c, certContext, certCtx)

	certDoc, err := issueCertificate(c, params)
	if err != nil {
		return nil, err
	}
	return certDoc.toModel(), nil
}
