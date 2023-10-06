package certtemplate

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func applyCertificateCapabilities(cap ns.NamespaceCertificateTemplateCapabilities, templateID models.Identifier,
	req models.CertificateTemplateParameters) (*CertificateTemplateDoc, error) {
	if templateID.HasReservedIDOrPrefix() && !cap.AllowedReservedNames.Contains(templateID) {
		return nil, fmt.Errorf("%w:reserved template id", common.ErrStatusBadRequest)
	}
	doc := CertificateTemplateDoc{}

	if cap.AllowedIssuerNamespaces.Size() == 1 {
		doc.IssuerNamespaceID = cap.AllowedIssuerNamespaces.Items()[0]
	} else {
		if req.Issuer == nil {
			return nil, fmt.Errorf("%w:issuer namespace is required", common.ErrStatusBadRequest)
		}
		issuerNsID := profile.GetResourceNsIDForProfile(kmsdoc.NewDocIdentifier(kmsdoc.DocKindDirectoryObject, req.Issuer.ProfileId))
		if !cap.AllowedIssuerNamespaces.Contains(issuerNsID) {
			return nil, fmt.Errorf("%w:invalid issuer namespace", common.ErrStatusBadRequest)
		}
	}
	if req.Issuer == nil || req.Issuer.TemplateId == nil {
		doc.IssuerTemplateID = kmsdoc.StringDocIdentifier(kmsdoc.DocKindCertificateTemplate, string(ns.CertTemplateNameDefault))
	} else {
		doc.IssuerTemplateID = kmsdoc.NewDocIdentifier(kmsdoc.DocKindCertificateTemplate, *req.Issuer.TemplateId)
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

	doc.KeyProperties.Kty = cap.DefaultKeyType
	doc.KeyProperties.Alg = cap.DefaultRsaAlgorithm
	switch doc.KeyProperties.Kty {
	case models.KeyTypeRSA:
		doc.KeyProperties.KeySize = utils.ToPtr(cap.DefaultKeySize)
		doc.KeyProperties.Alg = cap.DefaultRsaAlgorithm
	case models.KeyTypeEC:
		doc.KeyProperties.Crv = utils.ToPtr(cap.DefaultCrv)
	}
	if req.KeyProperties != nil {
		if req.KeyProperties.Kty == models.KeyTypeEC && cap.RestrictKeyTypeRsa {
			return nil, fmt.Errorf("%w:invalid key type", common.ErrStatusBadRequest)
		}
		switch req.KeyProperties.Kty {
		case models.KeyTypeRSA:
			doc.KeyProperties.Kty = models.KeyTypeRSA
			doc.KeyProperties.KeySize = req.KeyProperties.KeySize
			if doc.KeyProperties.KeySize == nil {
				doc.KeyProperties.KeySize = utils.ToPtr(cap.DefaultKeySize)
			}
			switch *doc.KeyProperties.KeySize {
			case 2048, 3072, 4096:
				// ok
			default:
				doc.KeyProperties.KeySize = utils.ToPtr(cap.DefaultKeySize)
			}
			if req.KeyProperties.Alg != nil {
				doc.KeyProperties.Alg = *req.KeyProperties.Alg
				switch doc.KeyProperties.Alg {
				case models.AlgRS256,
					models.AlgRS384,
					models.AlgRS512:
					// ok
				default:
					doc.KeyProperties.Alg = cap.DefaultRsaAlgorithm
				}
			}
		case models.KeyTypeEC:
			doc.KeyProperties.Kty = models.KeyTypeEC
			doc.KeyProperties.Crv = req.KeyProperties.Crv
			if doc.KeyProperties.Crv == nil {
				doc.KeyProperties.Crv = utils.ToPtr(cap.DefaultCrv)
			}
			switch *doc.KeyProperties.Crv {
			case models.CurveNameP256,
				models.CurveNameP384:
				// ok
			default:
				doc.KeyProperties.Crv = utils.ToPtr(cap.DefaultCrv)
			}
		}
	}
	if doc.KeyProperties.Kty == models.KeyTypeEC {
		switch *doc.KeyProperties.Crv {
		case models.CurveNameP256:
			doc.KeyProperties.Alg = models.AlgES256
		case models.CurveNameP384:
			doc.KeyProperties.Alg = models.AlgES384
		}
	}

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
	doc.LifetimeTrigger = req.LifetimeTrigger
	if doc.LifetimeTrigger == nil || (doc.LifetimeTrigger.DaysBeforeExpiry == nil &&
		doc.LifetimeTrigger.LifetimePercentage == nil) {
		// apply default
		doc.LifetimeTrigger = &models.CertificateLifetimeTrigger{
			LifetimePercentage: utils.ToPtr(int32(80)),
		}
	} else if doc.LifetimeTrigger.LifetimePercentage != nil {
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
		doc.LifetimeTrigger = nil
	}

	digestSb := strings.Builder{}
	digestSb.WriteString(doc.IssuerNamespaceID.String())
	digestSb.WriteString(doc.IssuerTemplateID.String())
	digestSb.WriteString(doc.SubjectCommonName)
	digestSb.WriteString(string(doc.KeyProperties.Alg))
	digestSb.WriteString(string(doc.KeyProperties.Kty))
	switch doc.KeyProperties.Kty {
	case models.KeyTypeRSA:
		digestSb.WriteString(fmt.Sprintf("%d", *doc.KeyProperties.KeySize))
	case models.KeyTypeEC:
		digestSb.WriteString(string(*doc.KeyProperties.Crv))
	}
	digest := md5.Sum([]byte(digestSb.String()))
	doc.Digest = digest[:]

	return &doc, nil
}

// PutCertificateTemplate implements CertificateTemplateService.
func validatePutRequest(c common.ServiceContext,
	templateID models.Identifier,
	req models.CertificateTemplateParameters) (*CertificateTemplateDoc, error) {

	if !templateID.IsValid() {
		return nil, fmt.Errorf("%w:invalid template id", common.ErrStatusBadRequest)
	}

	pcs := profile.GetProfileContextService(c)
	nsID := pcs.GetResourceDocNsID()

	caps, err := ns.GetNamespaceCapabilities(nsID)
	if err != nil {
		return nil, fmt.Errorf("%w:bad certificate requester", common.ErrStatusBadRequest)
	}
	certCaps := caps.GetAllowedCertificateIssuersForTemplate(templateID, pcs.GetRequestProfileType())
	doc, err := applyCertificateCapabilities(certCaps, templateID, req)
	if err != nil {
		return doc, err
	}

	profile, err := pcs.GetSelfProfileDoc(c)
	if err != nil {
		return nil, err
	}
	if profile.ProfileType != pcs.GetRequestProfileType() {
		return nil, fmt.Errorf("%w:invalid profile: type mismatch", common.ErrStatusBadRequest)
	}

	return doc, nil
}
