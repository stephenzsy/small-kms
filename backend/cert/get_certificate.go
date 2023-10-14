package cert

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	certtemplate "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func NewCertificateID(certId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindCert, certId)
}

func NewLatestCertificateForTemplateID(templateId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate, templateId)
}

func ReadCertDocByLocator(c context.Context, locator shared.ResourceLocator) (*CertDoc, error) {
	certDoc := &CertDoc{}
	err := kmsdoc.Read(c, locator, certDoc)
	return certDoc, err
}

func ApiGetCertificate(c RequestContext, certificateId shared.Identifier, params models.GetCertificateParams) error {
	cert, err := getCertificate(c, certificateId, params)
	if err != nil {
		return err
	}
	return c.JSON(200, cert)
}

func getLatestCertificateByTemplateDoc(c RequestContext, templateLocator shared.ResourceLocator) (doc *CertDoc, err error) {
	doc = &CertDoc{}
	err = kmsdoc.Read[*CertDoc](c,
		shared.NewResourceLocator(templateLocator.GetNamespaceID(), shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate, templateLocator.GetID().Identifier())), doc)
	return
}

func getCertificate(c RequestContext, certificateId shared.Identifier, params models.GetCertificateParams) (*shared.CertificateInfo, error) {
	var certDocLocator shared.ResourceLocator
	if certificateId.IsUUID() {
		nsID := ns.GetNamespaceContext(c).GetID()
		certDocLocator = shared.NewResourceLocator(nsID, NewCertificateID(certificateId))
	} else {
		return nil, fmt.Errorf("%w: invalid certificate ID: %s", common.ErrStatusBadRequest, certificateId)
	}

	certDoc, err := ReadCertDocByLocator(c, certDocLocator)
	if err != nil {
		return nil, err
	}
	m := certDoc.toModel()

	if params.IncludeCertificate != nil && (*params.IncludeCertificate == models.IncludePEM || *params.IncludeCertificate == models.IncludeJWK) {
		// fetch cert from blob
		pemBlob, err := certDoc.fetchCertificatePEMBlob(c)
		if err != nil {
			return m, err
		}
		m.Pem = utils.ToPtr(string(pemBlob))
		switch *params.IncludeCertificate {
		case models.IncludePEM:
		case models.IncludeJWK:
		}
		// attach blob
	}

	return m, nil
}

// This call has writes, please do not use for regular query
func GetAuthorizedLatestCertByTemplateID(c context.Context, templateID shared.Identifier) (*CertDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	if !templateID.IsUUID() || templateID.UUID().Version() != 5 {
		return ReadCertDocByLocator(c, shared.NewResourceLocator(nsID, NewLatestCertificateForTemplateID(templateID)))
	}

	// read linked doc
	localTemplateLocator := shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCertTemplate, templateID))
	localTemplateDoc, err := certtemplate.GetCertificateTemplateDoc(c, localTemplateLocator)
	if err != nil {
		return nil, err
	}
	if localTemplateDoc.Owner == nil {
		return nil, fmt.Errorf("%w: template is not linked", common.ErrStatusBadRequest)
	} else if localTemplateDoc.LinkProperties == nil || localTemplateDoc.LinkProperties.Usage != models.LinkedCertificateTemplateUsageClientAuthorization {
		return nil, fmt.Errorf("%w: template is not linked for client authorization", common.ErrStatusBadRequest)
	}
	remoteTemplateLocator := *localTemplateDoc.Owner
	remoteCertLinkLocator := remoteTemplateLocator.WithIDKind(shared.ResourceKindLatestCertForTemplate)
	remoteCertLinkDoc, err := ReadCertDocByLocator(c, remoteCertLinkLocator)
	if err != nil {
		return nil, err
	}
	// create link
	targetFinalLocator := remoteCertLinkDoc.GetLocator()
	targetCertDoc, err := ReadCertDocByLocator(c, targetFinalLocator)
	if err != nil {
		return nil, err
	}

	linkedCertDoc := *targetCertDoc
	linkedCertDoc.NamespaceID = nsID
	linkedCertDoc.ID = targetFinalLocator.GetID()
	linkedCertDoc.Owner = &targetFinalLocator
	linkedCertDoc.Template = localTemplateLocator

	eCtx := common.ElevateContext(c)
	err = kmsdoc.Upsert(eCtx, &linkedCertDoc)
	if err != nil {
		return nil, err
	}
	patchOps := azcosmos.PatchOperations{}
	if targetCertDoc.Owns == nil {
		patchOps.AppendSet(kmsdoc.PathPathOwns, map[shared.NamespaceIdentifier]shared.ResourceLocator{
			nsID: linkedCertDoc.GetLocator(),
		})
	} else {
		patchOps.AppendSet(fmt.Sprintf("%s/%s", kmsdoc.PathPathOwns, nsID), linkedCertDoc.GetLocator())
	}
	err = kmsdoc.Patch(eCtx, targetCertDoc, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: &targetCertDoc.ETag,
	})
	return &linkedCertDoc, err
}
