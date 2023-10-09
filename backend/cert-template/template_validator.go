package certtemplate

import (
	"crypto/md5"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func applyCertificateCapabilities(cap ns.NamespaceCertificateTemplateCapabilities, locator models.ResourceLocator,
	req models.CertificateTemplateParameters) (*CertificateTemplateDoc, error) {
	templateID := locator.GetID().Identifier()
	if templateID.HasReservedIDOrPrefix() {
		if _, contains := cap.AllowedReservedNames[templateID]; !contains {
			return nil, fmt.Errorf("%w:template id is reserved", common.ErrStatusBadRequest)
		}
	}
	doc := CertificateTemplateDoc{
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: locator.GetNamespaceID(),
			ID:          locator.GetID(),
		},
	}

	if req.IssuerTemplate != nil {
		if !cap.AllowedIssuerNamespaces.Contains(req.IssuerTemplate.GetNamespaceID()) {
			return nil, fmt.Errorf("%w:invalid issuer namespace", common.ErrStatusBadRequest)
		}
		reqResID := req.IssuerTemplate.GetID()
		if reqResID.Kind() != models.ResourceKindCertTemplate || !reqResID.Identifier().IsValid() {
			return nil, fmt.Errorf("%w:invalid issuer template", common.ErrStatusBadRequest)
		}
		doc.IssuerTemplate = *req.IssuerTemplate
	} else if cap.AllowedIssuerNamespaces.Size() == 1 {
		doc.IssuerTemplate = common.NewLocator(
			cap.AllowedIssuerNamespaces.Items()[0],
			common.NewIdentifierWithKind(models.ResourceKindCertTemplate, common.StringIdentifier(ns.CertTemplateNameDefault)))
	} else {
		return nil, fmt.Errorf("%w:issuer namespace is required", common.ErrStatusBadRequest)
	}
	if cap.SelfSigned {
		doc.IssuerTemplate = common.NewLocator(doc.NamespaceID, doc.ID)
	}

	if req.Usages == nil || len(req.Usages) == 0 {
		doc.Usages = cap.AllowedUsages.Items()
	} else {
		intersect := cap.AllowedUsages.Intersection(utils.NewSet[models.CertificateUsage](req.Usages...))
		if intersect.Size() == 0 {
			return nil, fmt.Errorf("%w:invalid certificate usages", common.ErrStatusBadRequest)
		}
		doc.Usages = intersect.Items()
	}

	if req.KeyProperties != nil && req.KeyProperties.Kty == models.KeyTypeEC && cap.RestrictKeyTypeRsa {
		return nil, fmt.Errorf("%w:invalid key type", common.ErrStatusBadRequest)
	}
	doc.KeySpec.initWithCreateTemplateInput(req.KeyProperties, CertKeySpec{
		Alg:     cap.DefaultRsaAlgorithm,
		Kty:     cap.DefaultKeyType,
		KeySize: &cap.DefaultKeySize,
		Crv:     &cap.DefaultCrv,
	})

	if cap.HasKeyStore {
		doc.KeyStorePath = req.KeyStorePath
		if doc.KeyStorePath == nil || *doc.KeyStorePath == "" {
			return nil, fmt.Errorf("%w:key store path is required", common.ErrStatusBadRequest)
		}
	}

	s, hasTemplate, err := preprocessTemplate(req.SubjectCommonName)
	if hasTemplate && !cap.AllowVariables {
		return nil, fmt.Errorf("%w:template variables are not allowed", common.ErrStatusBadRequest)
	}
	if err != nil {
		return nil, fmt.Errorf("%w:invalid subject common name", common.ErrStatusBadRequest)
	}
	if s == "" {
		return nil, fmt.Errorf("%w:subject common name is required", common.ErrStatusBadRequest)
	}
	doc.SubjectCommonName = s
	doc.ValidityInMonths = int32(cap.DefaultMaxValidityInMonths)
	if req.ValidityInMonths != nil && *req.ValidityInMonths != 0 {
		doc.ValidityInMonths = *req.ValidityInMonths
		if doc.ValidityInMonths < 0 {
			doc.ValidityInMonths = int32(1)
		} else if doc.ValidityInMonths > 120 {
			doc.ValidityInMonths = 120
		}
	}
	doc.LifetimeTrigger.DaysBeforeExpiry = nil
	doc.LifetimeTrigger.LifetimePercentage = utils.ToPtr(int32(80))
	if req.LifetimeTrigger != nil && (req.LifetimeTrigger.DaysBeforeExpiry != nil || req.LifetimeTrigger.LifetimePercentage != nil) {
		doc.LifetimeTrigger = *req.LifetimeTrigger
		if doc.LifetimeTrigger.LifetimePercentage != nil {
			if *doc.LifetimeTrigger.LifetimePercentage < 50 {
				doc.LifetimeTrigger.LifetimePercentage = utils.ToPtr(int32(50))
			}
			if *doc.LifetimeTrigger.LifetimePercentage >= 100 {
				doc.LifetimeTrigger.DaysBeforeExpiry = utils.ToPtr(int32(0))
				doc.LifetimeTrigger.LifetimePercentage = nil
			}
		} else if doc.LifetimeTrigger.DaysBeforeExpiry != nil {
			if *doc.LifetimeTrigger.DaysBeforeExpiry < 0 {
				doc.LifetimeTrigger.DaysBeforeExpiry = utils.ToPtr(int32(-1))
				doc.LifetimeTrigger.LifetimePercentage = nil
				// disabled
			} else if *doc.LifetimeTrigger.DaysBeforeExpiry > 15*doc.ValidityInMonths {
				doc.LifetimeTrigger.DaysBeforeExpiry = utils.ToPtr(int32(15 * doc.ValidityInMonths))
			}
		}
	}

	doc.Digest = doc.computeFieldsDigest()
	return &doc, nil
}

func (doc *CertificateTemplateDoc) computeFieldsDigest() []byte {
	digest := md5.New()
	digest.Write([]byte(doc.IssuerTemplate.String()))
	digest.Write([]byte(doc.SubjectCommonName))
	digest.Write([]byte(string(doc.KeySpec.Alg)))
	digest.Write([]byte(string(doc.KeySpec.Kty)))
	switch doc.KeySpec.Kty {
	case models.KeyTypeRSA:
		digest.Write([]byte(fmt.Sprintf("%d", *doc.KeySpec.KeySize)))
	case models.KeyTypeEC:
		digest.Write([]byte(string(*doc.KeySpec.Crv)))
	}
	return digest.Sum(nil)
}

// PutCertificateTemplate implements CertificateTemplateService.
func validatePutRequest(c RequestContext,
	locator models.ResourceLocator,
	req models.CertificateTemplateParameters) (*CertificateTemplateDoc, error) {

	nsID := locator.GetNamespaceID()
	certCaps := ns.GetAllowedCertificateIssuersForTemplate(locator)
	doc, err := applyCertificateCapabilities(certCaps, locator, req)
	if err != nil {
		return doc, err
	}

	profile, err := profile.GetResourceProfileDoc(c)
	if err != nil {
		return nil, err
	}
	if profile.ProfileType != nsID.Kind() {
		return nil, fmt.Errorf("%w:invalid profile: type mismatch", common.ErrStatusBadRequest)
	}

	return doc, nil
}
