package cert

import (
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

	certDoc, err := createCertificate(c, params)
	if err != nil {
		return nil, err
	}
	certDoc, err = issueCertificate(c, certDoc, params)
	if err != nil {
		return nil, err
	}
	return certDoc.toModel(), nil
}
