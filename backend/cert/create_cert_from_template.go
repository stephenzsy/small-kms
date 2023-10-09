package cert

import (
	"fmt"

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

func issueCertificate(c RequestContext,
	certDoc *CertDoc) (*CertDoc, error) {

	ctc := ct.GetCertificateTemplateContext(c)
	nsID := ns.GetNamespaceContext(c).GetID()
	// verify template
	if tmplDoc, err := ctc.GetCertificateTemplateDoc(c); err != nil {
		return nil, err
	} else if utils.IsTimeNotNilOrZero(tmplDoc.Deleted) {
		return nil, fmt.Errorf("%w: template not found", common.ErrStatusNotFound)
	}
	// verify issuer template still active
	if issuerTmplDoc, err := ct.GetCertificateTemplateDoc(c, certDoc.Issuer); err != nil {
		return nil, err
	} else if utils.IsTimeNotNilOrZero(issuerTmplDoc.Deleted) {
		return nil, fmt.Errorf("%w: issuer template not found", common.ErrStatusNotFound)
	}
	// verify profile
	if pdoc, err := profile.GetResourceProfileDoc(c); err != nil {
		return nil, err
	} else if utils.IsTimeNotNilOrZero(pdoc.Deleted) {
		return nil, fmt.Errorf("%w: profile not found", common.ErrStatusNotFound)
	}
	// verify graph
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot, models.NamespaceKindCaInt:
		// ok
	default:
		// verify graph
		gc, err := c.ServiceClientProvider().MsGraphDelegatedClient(c)
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

	var csrProvider CertificateRequestProvider
	var signerProvider SignerProvider
	var storageProvider StorageProvider = &azBlobStorageProvider{
		blobKey: fmt.Sprintf("%s/%s.pem", *certDoc.KeyStorePath, certDoc.ID.Identifier()),
	}

	switch nsID.Kind() {
	case models.NamespaceKindServicePrincipal:
		certDoc.CertSpec.keyExportable = true
	default:
		certDoc.CertSpec.keyExportable = false
	}

	switch nsID.Kind() {
	case models.NamespaceKindCaRoot:
		if certDoc.Issuer != certDoc.Template {
			return nil, fmt.Errorf("invalid issuer template for root ca, must be self")
		}
		selfSignProvider := newAzKeysSelfSignerProvider(certDoc)
		signerProvider = selfSignProvider
		csrProvider = selfSignProvider
	case models.NamespaceKindCaInt,
		models.NamespaceKindServicePrincipal:
		issuerDoc, err := certDoc.readIssuerCertDoc(c)
		if err != nil {
			return nil, err
		}
		csrProvider = newAzCertsCsrProvider(certDoc)
		signerProvider = newAzKeysExistingCertSigner(issuerDoc)
	default:
		return nil, fmt.Errorf("%w: invalid namespace kind", common.ErrStatusBadRequest)
	}

	patch, err := signCertificate(c, csrProvider, signerProvider, storageProvider)
	if err != nil {
		return nil, err
	}
	err = certDoc.patchSigned(c, patch)
	if err != nil {
		return nil, err
	}
	certDocLatestLinkLocator := models.NewResourceLocator(nsID, models.NewResourceIdentifier(models.ResourceKindLatestCertForTemplate,
		certDoc.Template.GetID().Identifier()))

	_, err = kmsdoc.UpsertAliasWithSnapshot(c, certDoc, certDocLatestLinkLocator)
	if err != nil {
		return nil, err
	}

	return certDoc, nil
}
