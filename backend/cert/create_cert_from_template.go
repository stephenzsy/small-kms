package cert

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func shouldCreateNewCertificate(c RequestContext, templateDoc *ct.CertificateTemplateDoc) (*CertDoc, bool, error) {
	certDoc, err := getLatestCertificateByTemplateDoc(c, templateDoc.GetLocator())
	if err != nil {
		if !errors.Is(err, common.ErrStatusNotFound) {
			return nil, false, err
		}
		return nil, true, nil
	}

	if certDoc.Deleted != nil {
		log.Info().Msg("should create new cert: latest certificate is deleted")
		return certDoc, true, nil
	}
	if certDoc.NotAfter.Time().Before(time.Now()) {
		log.Info().Msg("should create new cert: latest certificate is expired")
		return certDoc, true, nil
	}
	if !slices.Equal(certDoc.TemplateDigest, templateDoc.Digest) {
		log.Info().Msg("should create new cert: template digest mismatch")
		return certDoc, true, nil
	}

	if templateDoc.LifetimeTrigger.DaysBeforeExpiry != nil {
		cutOff := certDoc.NotAfter.Time().AddDate(0, 0, 0-int(*templateDoc.LifetimeTrigger.DaysBeforeExpiry))
		shouldRenew := cutOff.Before(time.Now())
		log.Info().Msgf("should create new cert: eval days before expiry: %v", shouldRenew)
		return certDoc, shouldRenew, nil
	}

	if templateDoc.LifetimeTrigger.LifetimePercentage != nil {
		duration := certDoc.NotAfter.Time().Sub(certDoc.NotBefore.Time())
		cutOff := certDoc.NotBefore.Time().Add(duration * time.Duration(*templateDoc.LifetimeTrigger.LifetimePercentage) / 100)
		shouldRenew := cutOff.Before(time.Now())
		log.Info().Msgf("should create new cert: eval lifetime percentage: %v", shouldRenew)
		return certDoc, shouldRenew, nil
	}

	log.Info().Msg("should not create new cert: current is within range")
	return certDoc, false, nil
}

func issueCertificate(c RequestContext,
	certDoc *CertDoc) (context.Context, *CertDoc, error) {

	bad := func(e error) (context.Context, *CertDoc, error) {
		return nil, nil, e
	}

	nsID := ns.GetNamespaceContext(c).GetID()

	// verify issuer template still active
	if issuerTmplDoc, err := ct.GetCertificateTemplateDoc(c, certDoc.Issuer); err != nil {
		return bad(err)
	} else if utils.IsTimeNotNilOrZero(issuerTmplDoc.Deleted) {
		return bad(fmt.Errorf("%w: issuer template not found", common.ErrStatusNotFound))
	}
	// verify profile
	if pdoc, err := profile.GetResourceProfileDoc(c); err != nil {
		return bad(err)
	} else if utils.IsTimeNotNilOrZero(pdoc.Deleted) {
		return bad(fmt.Errorf("%w: profile not found", common.ErrStatusNotFound))
	}
	// verify graph
	switch nsID.Kind() {
	case shared.NamespaceKindCaRoot, shared.NamespaceKindCaInt:
		// ok
	default:
		// verify graph
		gc, err := common.GetAdminServerRequestClientProvider(c).MsGraphClient()
		if err != nil {
			return bad(err)
		}
		dirObj, err := gc.DirectoryObjects().ByDirectoryObjectId(nsID.Identifier().String()).Get(c, nil)
		if err != nil {
			return bad(err)
		}
		pdoc, err := profile.StoreProfile(c, dirObj, nil, nil)
		if err != nil {
			return bad(err)
		}
		if pdoc.ProfileType != nsID.Kind() {
			return bad(fmt.Errorf("%w: invalid profile type, mismatch", common.ErrStatusBadRequest))
		}
	}

	var csrProvider CertificateRequestProvider
	var signerProvider SignerProvider
	var storageProvider StorageProvider = &azBlobStorageProvider{
		blobKey: fmt.Sprintf("%s/%s.pem", *certDoc.KeyStorePath, certDoc.ID.Identifier()),
	}

	switch nsID.Kind() {
	case shared.NamespaceKindServicePrincipal:
		certDoc.CertSpec.keyExportable = true
	default:
		certDoc.CertSpec.keyExportable = false
	}

	switch nsID.Kind() {
	case shared.NamespaceKindCaRoot:
		if certDoc.Issuer != certDoc.Template {
			return bad(fmt.Errorf("invalid issuer template for root ca, must be self"))
		}
		selfSignProvider := newAzKeysSelfSignerProvider(certDoc)
		signerProvider = selfSignProvider
		csrProvider = selfSignProvider
	case shared.NamespaceKindCaInt,
		shared.NamespaceKindServicePrincipal:
		issuerDoc, err := certDoc.readIssuerCertDoc(c)
		if err != nil {
			return bad(err)
		}
		csrProvider = newAzCertsCsrProvider(certDoc, false)
		signerProvider = newAzKeysExistingCertSigner(issuerDoc)
	default:
		return bad(fmt.Errorf("%w: invalid namespace kind", common.ErrStatusBadRequest))
	}
	{
		c := ctx.Elevate(c)
		var patch *CertDocSigningPatch
		patch, err := signCertificate(c, csrProvider, signerProvider, storageProvider)
		if err != nil {
			return bad(err)
		}
		if patchOps := certDoc.patchSigned(c, patch); patchOps == nil {
			// persist doc for new certs
			err = kmsdoc.Create(c, certDoc)
		} else {
			// patch doc
			err = kmsdoc.Patch(c, certDoc, *patchOps, nil)
		}
		if err != nil {
			// patch failed, but cert is signed
			log.Error().Err(err).Msgf("Cert is issued, but failed to create/patch cert doc: %s", certDoc.GetLocator())
			// TODO: disable cert on key vault
			return bad(err)
		}

		// create template link in this block for now
		certDocLatestLinkLocator := shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate,
			certDoc.Template.GetID().Identifier()))

		_, err = kmsdoc.UpsertAliasWithSnapshot(c, certDoc, certDocLatestLinkLocator)
		if err != nil {
			return bad(err)
		}

		return c, certDoc, nil
	}
}
