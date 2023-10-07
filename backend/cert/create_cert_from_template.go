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
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
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
		Status:            CertStatusInitialized,
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
	certDoc *CertDoc,
	params models.IssueCertificateFromTemplateParams) (*CertDoc, error) {

	ctc := ct.GetCertificateTemplateContext(c)
	nsID := ns.GetNamespaceContext(c).GetID()
	// verify template
	if tmplDoc, err := ctc.GetCertificateTemplateDoc(c); err != nil {
		return nil, err
	} else if utils.IsTimeNotNilOrZero(tmplDoc.Deleted.TimePtr()) {
		return nil, fmt.Errorf("%w: template not found", common.ErrStatusNotFound)
	}
	// verify issuer template still active
	if issuerTmplDoc, err := ct.GetCertificateTemplateDoc(c, certDoc.Issuer); err != nil {
		return nil, err
	} else if utils.IsTimeNotNilOrZero(issuerTmplDoc.Deleted.TimePtr()) {
		return nil, fmt.Errorf("%w: issuer template not found", common.ErrStatusNotFound)
	}
	// verify profile
	if pdoc, err := profile.GetResourceProfileDoc(c); err != nil {
		return nil, err
	} else if utils.IsTimeNotNilOrZero(pdoc.Deleted.TimePtr()) {
		return nil, fmt.Errorf("%w: profile not found", common.ErrStatusNotFound)
	}
	// verify graph
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot, models.NamespaceKindCaInt:
		// ok
	default:
		// verify graph
		gc, err := common.GetClientProvider(c).MsGraphDelegatedClient(c)
		if err != nil {
			return nil, err
		}
		dirObj, err := gc.DirectoryObjects().ByDirectoryObjectId(nsID.Identifier().String()).Get(c, nil)
		if err != nil {
			return nil, err
		}
		pdoc, err := profile.StoreProfile(c, dirObj)
		if err != nil {
			return nil, err
		}
		if pdoc.ProfileType != nsID.Kind() {
			return nil, fmt.Errorf("%w: invalid profile type, mismatch", common.ErrStatusBadRequest)
		}
	}

	// now need to load signer certificate
	if nsID.Kind() == models.NamespaceKindCaRoot {
		if certDoc.Issuer != certDoc.Template {
			return nil, fmt.Errorf("invalid issuer template for root ca, must be self")
		}
		// root is self signed, no need to load signer certificate

	}

	_, err := signCertificate(c, nil, nil, certDoc, nil)
	if err != nil {
		return nil, err
	}

	panic("unimplemented")
}
