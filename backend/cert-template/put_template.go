package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

// PutCertificateTemplate implements CertificateTemplateService.
func PutCertificateTemplate(c RequestContext,
	req models.CertificateTemplateParameters) (*models.CertificateTemplateComposed, error) {

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	locator := GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c)
	doc, err := validatePutRequest(c, locator, req)
	if err != nil {
		return nil, err
	}

	doc.SchemaVersion = 1
	err = kmsdoc.Upsert(c, doc)
	if err != nil {
		return nil, err
	}

	return doc.toModel(), nil
}
