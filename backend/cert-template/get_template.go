package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/profile"
)

func getCertificateTemplateDoc(c common.ServiceContext,
	templateID models.Identifier) (doc *CertificateTemplateDoc, err error) {
	pc := profile.GetProfileContext(c)
	nsID := pc.GetResourceDocNsID()

	doc = new(CertificateTemplateDoc)
	err = kmsdoc.Read(c, nsID, kmsdoc.NewDocIdentifier(kmsdoc.DocKindCertificateTemplate, templateID), doc)
	return
}

// PutCertificateTemplate implements CertificateTemplateService.
func (s *certTmplService) GetCertificateTemplate(c common.ServiceContext,
	templateID models.Identifier) (*models.CertificateTemplate, error) {

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	doc, err := getCertificateTemplateDoc(c, templateID)
	if err != nil {
		return nil, err
	}

	return doc.toModel(), nil
}
