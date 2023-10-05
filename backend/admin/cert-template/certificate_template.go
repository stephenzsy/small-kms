package certtemplate

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/admin/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertificateTemplate interface {
	IsEnabled() bool

	CerateCert() (cert.Certificate, error)
	CreateCertWithVariables(TemplateVarData) (cert.Certificate, error)
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
func (*certificateTemplate) CreateCertWithVariables(TemplateVarData) (cert.Certificate, error) {
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

type CreateCertificateTemplateParameters struct {
	NamespaceID             uuid.UUID
	TemplateID              uuid.UUID
	Features                utils.Set[CertificateTemplateFlag]
	DisplayName             string
	IssuerNamespaceID       uuid.UUID
	IssuerTemplateID        uuid.UUID
	KeyProperties           CertificateTemplateDocKeyProperties
	KeyStorePath            string
	Subject                 CertificateTemplateDocSubject
	SubjectAlternativeNames *CertificateTemplateDocSANs
	ValidityInMonths        int32
	LifetimeTrigger         *CertificateTemplateDocLifeTimeTrigger
}

func CreateTemplate(params CreateCertificateTemplateParameters) (CertificateTemplate, error) {
	panic("unimplemented")
}