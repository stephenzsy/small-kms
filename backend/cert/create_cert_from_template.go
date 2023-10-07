package cert

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

var (
	ErrInvalidContext = fmt.Errorf("invalid context")
)

func createCertificate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*CertDoc, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	ctc := ct.GetCertificateTemplateContext(c)
	tmpl, err := ctc.GetCertificateTemplateDoc(c)
	if err != nil {
		return nil, err
	}

	certID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	doc := CertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: nsID,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCert, common.UUIDIdentifier(certID)),
		},
		Status:            CertStatusInitial,
		SerialNumber:      SerialNumberStorable(certID[:]),
		SubjectCommonName: tmpl.SubjectCommonName,
		Usages:            tmpl.Usages,
		KeySpec:           tmpl.KeySpec,
		KeyStorePath:      tmpl.KeyStorePath,
		Template:          tmpl.GetLocator(),
		Issuer:            tmpl.IssuerTemplate,
		NotBefore:         kmsdoc.TimeStorable(now),
		NotAfter:          kmsdoc.TimeStorable(now.AddDate(0, int(tmpl.ValidityInMonths), 0)),
	}

	return &doc, nil
}

func issueCertificate(c common.ServiceContext,
	params models.IssueCertificateFromTemplateParams) (*CertDoc, error) {

	_ = ct.GetCertificateTemplateContext(c)

	panic("unimplemented")
}
