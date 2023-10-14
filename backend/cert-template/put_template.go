package certtemplate

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/shared"
)

// PutCertificateTemplate implements CertificateTemplateService.
func PutCertificateTemplate(c RequestContext,
	req models.CertificateTemplateParameters) (*models.CertificateTemplateComposed, error) {

	locator := GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c)

	if locator.GetID().Kind() != shared.ResourceKindCertTemplate {
		return nil, fmt.Errorf("%w:invalid resource type: %s, expected: %s", common.ErrStatusBadRequest, locator.GetID().Kind(), shared.ResourceKindCertTemplate)
	}
	if locator.GetID().Identifier().IsUUID() && locator.GetID().Identifier().UUID().Version() == 5 {
		return nil, fmt.Errorf("%w:invalid resource ID: %s", common.ErrStatusBadRequest, locator.GetID().Identifier().String())
	}

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
