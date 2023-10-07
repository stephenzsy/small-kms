package certtemplate

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

func getCertificateTemplateDocLocator(nsID models.NamespaceID, templateID common.Identifier) models.ResourceLocator {
	return common.NewLocator(nsID, common.NewIdentifierWithKind(models.ResourceKindCertTemplate, templateID))
}

func GetCertificateTemplateDoc(c common.ServiceContext,
	locator models.ResourceLocator) (doc *CertificateTemplateDoc, err error) {

	if locator.GetID().Kind() != models.ResourceKindCertTemplate {
		return nil, fmt.Errorf("invalid resource type: %s, expected: %s", locator.GetID().Kind(), models.ResourceKindCertTemplate)
	}

	doc = new(CertificateTemplateDoc)
	err = kmsdoc.Read(c, locator, doc)
	return
}

// PutCertificateTemplate implements CertificateTemplateService.
func GetCertificateTemplate(c common.ServiceContext,
) (*models.CertificateTemplateComposed, error) {

	templateLocator := GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c)
	doc, err := GetCertificateTemplateDoc(c, templateLocator)
	if err != nil {
		return nil, err
	}

	return doc.toModel(), nil
}
