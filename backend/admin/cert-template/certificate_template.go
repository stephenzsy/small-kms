package certtemplate

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/admin/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertificateTemplate interface {
	IsEnabled() bool

	CerateCert() (cert.Certificate, error)
	CreateCertWithVariable() (cert.Certificate, error)
}

type certificateTemplate struct {
	doc *CertificateTemplateDoc
}

// IsEnabled implements CertificateTemplate.
func (t *certificateTemplate) IsEnabled() bool {
	return t.doc.Deleted == nil || t.doc.Deleted.IsZero()
}

// CerateCert implements CertTmpl.
func (*certificateTemplate) CerateCert() (cert.Certificate, error) {
	panic("unimplemented")
}

// CreateCertWithVariable implements CertTmpl.
func (*certificateTemplate) CreateCertWithVariable() (cert.Certificate, error) {
	panic("unimplemented")
}

func LoadCertifictateTemplate(c common.ServiceContext, nsID uuid.UUID, templateID uuid.UUID) (CertificateTemplate, error) {
	bad := func(e error) (CertificateTemplate, error) {
		return nil, e
	}

	doc := CertificateTemplateDoc{}
	if err := kmsdoc.AzCosmosReadItem(c, nsID, kmsdoc.NewKmsDocID(kmsdoc.DocTypeCertTemplate, templateID), &doc); err != nil {
		return bad(err)
	}

	return &certificateTemplate{doc: &doc}, nil
}
