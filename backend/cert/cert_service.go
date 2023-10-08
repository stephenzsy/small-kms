package cert

import (
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

/*
type certServiceContextKey string

const (
	certContext certServiceContextKey = "certContext"
)
*/
// IssueCertificateFromTemplate implements CertificateService.
func IssueCertificateFromTemplate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*models.CertificateInfoComposed, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	ctc := ct.GetCertificateTemplateContext(c)
	tmpl, err := ctc.GetCertificateTemplateDoc(c)
	if err != nil {
		return nil, err
	}

	certDoc, err := createCertificateDoc(nsID, tmpl, params)
	if err != nil {
		return nil, err
	}

	// persist document
	err = kmsdoc.Create(c, certDoc)
	if err != nil {
		return nil, err
	}

	certDoc, err = issueCertificate(c, certDoc, params)
	if err != nil {
		return nil, err
	}
	return certDoc.toModel(), nil
}
