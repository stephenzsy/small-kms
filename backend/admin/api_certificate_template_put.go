package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ct "github.com/stephenzsy/small-kms/backend/admin/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var (
	ErrCertificateTemplateVariable = errors.New("certificate template variable field is invalid")
)

func (s *adminServer) PutCertificateTemplateV2(c *gin.Context, namespaceId uuid.UUID, templateId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	isValid, needGraphValidation := isGraphValidationNeeded(namespaceId)
	var odataType graph.MsGraphOdataType
	if !isValid {
		respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace ID is not valid: %s", namespaceId))
		return
	}
	if needGraphValidation {
		// will check if directory object is already sync, sync will performed prior to issuing certificates
		graphObj, err := s.graphService.GetGraphProfileDoc(c, namespaceId, graph.MsGraphOdataTypeAny)
		if err != nil {
			common.RespondError(c, err)
			return
		}
		odataType = graphObj.GetOdataType()
	}

	templateParams := CertificateTemplateParameters{}
	err := c.Bind(&templateParams)
	if err != nil {
		respondPublicError(c, http.StatusBadRequest, err)
		return
	}

	doc, err := templateParams.validateAndToDoc(odataType, namespaceId, templateId)
	if err != nil {
		respondPublicError(c, http.StatusBadRequest, err)
		return
	}
	if doc == nil {
		respondPublicErrorMsg(c, http.StatusBadRequest, "no valid input")
		return
	}

	if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
		respondInternalError(c, err, fmt.Sprintf("failed to upsert certificate template in cosmos: %s", templateId))
		return
	}

	c.JSON(http.StatusOK, doc.toCertificateTemplate())
}

func (p *CertificateTemplateParameters) populateDocIssuer(doc *CertificateTemplateDoc, issuerNsType NamespaceTypeShortName) {
	doc.IssuerNamespaceID = p.Issuer.NamespaceID
	doc.IssuerNameSpaceType = issuerNsType
	if p.Issuer.TemplateID != nil {
		doc.IssuerTemplateID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeCertTemplate, *p.Issuer.TemplateID)
	} else {
		doc.IssuerTemplateID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeCertTemplate, uuid.Nil)
	}
}

func validateTemplateIdentifiers(nsID uuid.UUID, templateID uuid.UUID, name string, allowNamespaceDefault, allowEntraClientCredsDefault bool) (string, bool) {
	if templateID.Version() == 4 {
		// only allow non default prefixed name for user specified template ID
		return name, (name != "" && !strings.HasPrefix(name, "default"))
	} else if templateID.Version() == 5 {
		if allowNamespaceDefault && templateID == common.GetCanonicalCertificateTemplateID(nsID, common.DefaultCertTemplateName_GlobalDefault) {
			return string(common.DefaultCertTemplateName_GlobalDefault), true
		}
		if allowEntraClientCredsDefault && templateID == common.GetCanonicalCertificateTemplateID(nsID, common.DefaultCertTemplateName_ServicePrincipalClientCredential) {
			return string(common.DefaultCertTemplateName_ServicePrincipalClientCredential), true
		}
	}
	return "invalid", false
}

func (p *CertificateTemplateParameters) getTemplateFlags(odataType graph.MsGraphOdataType, nsID uuid.UUID, templateId uuid.UUID) (flags utils.Set[ct.CertificateTemplateFlag]) {
	return
}

func (p *CertificateTemplateParameters) toCreateCertificateTemplateParameters(odataType graph.MsGraphOdataType, nsID uuid.UUID, templateId uuid.UUID) (*ct.CreateCertificateTemplateParameters, error) {
	if p == nil {
		return nil, nil
	}

	tmplFlags := utils.NewSet[ct.CertificateTemplateFlag]()
	allowNamespace := false
	allowNamespaceDefault := false
	allowEntraClientCredsDefault := false
	var issuerNamespaceID uuid.UUID
	if isAllowedCaNamespace(nsID) {
		//tmplFlags.Add(ct.Flag)
		if isAllowedRootCaNamespace(nsID) {
			//	tmplFlags.Add(ct.CertificateTemplateFlagCA)
		}
		if isTestCA(nsID) {
			tmplFlags.Add(ct.CertTmplFlagTest)
			issuerNamespaceID = common.WellKnownID_TestRootCA
		} else {
			issuerNamespaceID = common.WellKnownID_RootCA
		}
		allowNamespace = true
		allowNamespaceDefault = true

		tmplFlags.Add(ct.CertTmplFlagHasKeyStore)
	} else {
		if !isAllowedIntCaNamespace(p.Issuer.NamespaceID) {
			return nil, fmt.Errorf("issuer namespace ID %s is not allowed to issue certificate", p.Issuer.NamespaceID)
		}
		issuerNamespaceID = p.Issuer.NamespaceID
		// usageFeatures := utils.NewSet[ct.CertificateTemplateFlag](
		// 	ct.CertTmplFlagDelegate,
		// )
		// switch p.Usage {
		// case UsageServerOnly:
		// 	usageFeatures.Add(ct.TemplateCertificateFeatureServer)
		// 	usageFeatures.Add(ct.CertTmplFlagHasKeyStore)
		// 	usageFeatures.Add(ct.CertTmplFlagKeyExportable)
		// case UsageClientOnly:
		// 	usageFeatures.Add(ct.TemplateCertificateFeatureClient)
		// case UsageAADClientCredential:
		// 	usageFeatures.Add(ct.CertTmplFlagRestrictKtyRsa)
		// 	usageFeatures.Add(ct.TemplateCertificateFeatureServer)
		// 	usageFeatures.Add(ct.TemplateCertificateFeatureClient)
		// default:
		// 	usageFeatures.Add(ct.CertTmplFlagHasKeyStore)
		// 	usageFeatures.Add(ct.CertTmplFlagKeyExportable)
		// 	usageFeatures.Add(ct.TemplateCertificateFeatureServer)
		// 	usageFeatures.Add(ct.TemplateCertificateFeatureClient)
		// }

		// namespaceFeatures := utils.NewSet[ct.CertificateTemplateFlag]()
		// switch odataType {
		// case graph.MsGraphOdataTypeUser:
		// 	namespaceFeatures.Add(ct.TemplateCertificateFeatureClient)
		// case graph.MsGraphOdataTypeGroup:
		// 	namespaceFeatures.Add(ct.CertTmplFlagDelegate)

		// 	namespaceFeatures.Add(ct.TemplateCertificateFeatureServer)
		// 	namespaceFeatures.Add(ct.TemplateCertificateFeatureClient)
		// 	allowNamespace = true
		// 	allowEntraClientCredsDefault = true
		// case graph.MsGraphOdataTypeServicePrincipal:
		// 	namespaceFeatures.Add(ct.TemplateCertificateFeatureServer)
		// 	namespaceFeatures.Add(ct.TemplateCertificateFeatureClient)

		// 	namespaceFeatures.Add(ct.CertTmplFlagRestrictKtyRsa)

		// 	namespaceFeatures.Add(ct.CertTmplFlagHasKeyStore)
		// 	namespaceFeatures.Add(ct.CertTmplFlagKeyExportable)

		// 	allowNamespace = true
		// 	allowEntraClientCredsDefault = true
		// 	allowNamespaceDefault = true
		// }
		// tmplFlags = usageFeatures.Intersection(namespaceFeatures)
	}

	if tmplFlags.Size() == 0 {
		//return nil, fmt.Errorf("certificate usage %s is not allowed", p.Usage)
	}

	if !allowNamespace {
		return nil, fmt.Errorf("namespace type %s is not allowed to create certificate template", odataType)
	}

	// display name
	displayName := p.DisplayName
	if fixedName, ok := validateTemplateIdentifiers(nsID, templateId, displayName, allowNamespaceDefault, allowEntraClientCredsDefault); ok {
		displayName = fixedName
	} else {
		return nil, fmt.Errorf("template ID %s is not valid for namespace: ", templateId, nsID)
	}

	// issuer template ID
	issuerTemplateId := common.GetCanonicalCertificateTemplateID(issuerNamespaceID, common.DefaultCertTemplateName_GlobalDefault)
	if p.Issuer.TemplateID != nil && *p.Issuer.TemplateID != uuid.Nil {
		issuerTemplateId = *p.Issuer.TemplateID
	}

	// resolve key properties
	keyProperties := ct.CertificateTemplateDocKeyProperties{}

	allowedKty, allowedKeySize, allowedCrv := utils.NewSet[models.JwtKty](models.KeyTypeRSA, models.KeyTypeEC),
		utils.NewSet[int](2048, 3072, 4096),
		utils.NewSet[models.JwtCrv](models.CurveNameP256, models.CurveNameP384)
	defaultKty, defaultKeySize := models.KeyTypeRSA, 2048
	if tmplFlags.Contains(ct.CertTmplFlagTest) {
		defaultKty = models.KeyTypeEC
	}
	switch {
	// case tmplFlags.Contains(ct.TemplateCertificateFeatureCA):
	// 	defaultKeySize = 4096
	case tmplFlags.Contains(ct.CertTmplFlagRestrictKtyRsa):
		allowedKty.Remove(models.KeyTypeEC)
		defaultKty = models.KeyTypeRSA
	default:
		defaultKty = models.KeyTypeRSA
		defaultKeySize = 2048
	}
	if p.KeyProperties == nil {
		keyProperties.Kty = defaultKty
		switch defaultKty {
		case models.KeyTypeRSA:
			keyProperties.KeySize = &defaultKeySize
		case models.KeyTypeEC:
			keyProperties.Crv = utils.ToPtr(models.CurveNameP384)
		}
	} else {
		if !allowedKty.Contains(p.KeyProperties.Kty) {
			return nil, fmt.Errorf("key type %s is not allowed", p.KeyProperties.Kty)
		}
		keyProperties.Kty = p.KeyProperties.Kty
		switch p.KeyProperties.Kty {
		case models.KeyTypeRSA:
			if p.KeyProperties.KeySize == nil || *p.KeyProperties.KeySize == 0 {
				keyProperties.KeySize = &defaultKeySize
			} else if !allowedKeySize.Contains(*p.KeyProperties.KeySize) {
				return nil, fmt.Errorf("key size %d is not allowed", *p.KeyProperties.KeySize)
			} else {
				keyProperties.KeySize = p.KeyProperties.KeySize
			}
		case models.KeyTypeEC:
			if p.KeyProperties.Crv != nil && !allowedCrv.Contains(*p.KeyProperties.Crv) {
				return nil, fmt.Errorf("curve %s is not allowed", *p.KeyProperties.Crv)
			}
			if p.KeyProperties.Crv == nil {
				keyProperties.Crv = utils.ToPtr(models.CurveNameP384)
			} else {
				keyProperties.Crv = p.KeyProperties.Crv
			}
		}
	}

	var keyStorePath string
	if tmplFlags.Contains(ct.CertTmplFlagHasKeyStore) && (p.KeyStorePath == nil || len(strings.TrimSpace(*p.KeyStorePath)) == 0) {
		return nil, fmt.Errorf("key store path must be specified")
	} else {
		keyStorePath = strings.TrimSpace(*p.KeyStorePath)
	}

	return &ct.CreateCertificateTemplateParameters{
		NamespaceID:       nsID,
		TemplateID:        templateId,
		Features:          tmplFlags,
		DisplayName:       displayName,
		IssuerNamespaceID: issuerNamespaceID,
		IssuerTemplateID:  issuerTemplateId,
		KeyProperties:     keyProperties,
		KeyStorePath:      keyStorePath,
	}, nil
}

func (p *CertificateTemplateParameters) validateAndToDoc(odataType graph.MsGraphOdataType, nsID uuid.UUID, templateId uuid.UUID) (*CertificateTemplateDoc, error) {
	if p == nil {
		return nil, nil
	}

	nsType := NSTypeAny
	switch {
	case isAllowedRootCaNamespace(nsID):
		nsType = NSTypeRootCA
	case isAllowedIntCaNamespace(nsID):
		nsType = NSTypeIntCA
	default:
		nsType = OdataTypeToNSType(odataType)
	}

	displayName := p.DisplayName
	// if fixedName, ok := validateTemplateIdentifiersOld(nsType, nsID, templateId, displayName); ok {
	// 	displayName = fixedName
	// } else {
	// 	return nil, fmt.Errorf("template ID %s is not valid for namespace type %s", templateId, nsType)
	// }

	doc := new(CertificateTemplateDoc)
	doc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeCertTemplate, templateId)
	doc.NamespaceID = nsID
	doc.DisplayName = displayName

	// validate and populate issuer, usage
	switch nsType {
	case NSTypeRootCA:
		if p.Issuer.NamespaceID != nsID {
			return nil, fmt.Errorf("root ca issuer namespace ID %s does not match namespace ID %s", p.Issuer.NamespaceID, nsID)
		}
		doc.IssuerNamespaceID = nsID
		if p.Issuer.TemplateID != nil && *p.Issuer.TemplateID != templateId {
			return nil, fmt.Errorf("root ca issuer template ID %s does not match template ID %s", *p.Issuer.TemplateID, templateId)
		}
		doc.IssuerTemplateID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeCertTemplate, templateId)
		doc.IssuerNameSpaceType = NSTypeRootCA
		// if p.Usage != UsageRootCA {
		// 	return nil, errors.New("root ca must be used for root ca certificate")
		// }
		doc.Usage = UsageRootCA
		doc.ValidityInMonths = 120 // default 10 years
	case NSTypeIntCA:
		if !isAllowedRootCaNamespace(p.Issuer.NamespaceID) {
			return nil, fmt.Errorf("intermediate ca issuer namespace ID %s is not a root ca namespace ID", p.Issuer.NamespaceID)
		}
		if isTestCA(nsID) && !isTestCA(p.Issuer.NamespaceID) {
			return nil, fmt.Errorf("test ca namespace ID %s can only issue certificates to test intermediate ca namespace", nsID)
		} else if !isTestCA(nsID) && isTestCA(p.Issuer.NamespaceID) {
			return nil, fmt.Errorf("production ca namespace ID %s can only issue certificates to production intermediate ca namespace", nsID)
		}
		p.populateDocIssuer(doc, NSTypeRootCA)
		// if p.Usage != UsageIntCA {
		// 	return nil, errors.New("intermediate ca must be used for intermediate ca certificate")
		// }
		doc.Usage = UsageIntCA
		doc.ValidityInMonths = 36 // default 3 years

	default:
		if !isAllowedIntCaNamespace(p.Issuer.NamespaceID) {
			return nil, fmt.Errorf("service principal/group issuer namespace ID %s is not an intermediate ca namespace ID", p.Issuer.NamespaceID)
		}
		p.populateDocIssuer(doc, NSTypeIntCA)
		// if p.Usage == UsageRootCA || p.Usage == UsageIntCA {
		// 	return nil, errors.New("root ca/intermediate ca must not be used for other certificate")
		// }
		// if p.Usage == UsageAADClientCredential && nsType != NSTypeServicePrincipal {
		// 	return nil, errors.New("AAD client credential certificate must be used for service principal")
		// }
		//doc.Usage = p.Usage
		doc.ValidityInMonths = 12 // default 1 year
	}

	// validate and populate key properties
	doc.KeyProperties.setDefault()
	//doc.KeyProperties.ReuseKey = p.ReuseKey
	switch nsType {
	case NSTypeRootCA,
		NSTypeIntCA:
		if isTestCA(nsID) {
			doc.KeyProperties.setECDSA(models.CurveNameP384)
		} else {
			doc.KeyProperties.setRSA(models.AlgRS384, 4096)
		}
		// ignore input key properties for CA
	default:
		if err := doc.KeyProperties.fromJwkProperties(p.KeyProperties); err != nil {
			return nil, err
		}
		// if p.Usage == UsageAADClientCredential && doc.KeyProperties.Kty == models.KeyTypeEC {
		// 	return nil, fmt.Errorf("AAD client credential certificate must use RSA key")
		// }
	}

	switch nsType {
	case NSTypeGroup:
		if p.KeyStorePath != nil && len(*p.KeyStorePath) > 0 {
			return nil, fmt.Errorf("group certificate must not specify key store path")
		}
	default:
		if p.KeyStorePath == nil || len(*p.KeyStorePath) == 0 {
			return nil, fmt.Errorf("key store path must be specified")
		}
		doc.KeyStorePath = p.KeyStorePath
	}

	// doc.Subject = CertificateTemplateDocSubject{CertificateSubject: p.Subject}
	// doc.SubjectAlternativeNames = sanitizeSANs(p.SubjectAlternativeNames)

	if p.ValidityInMonths != nil {
		if *p.ValidityInMonths < 0 && *p.ValidityInMonths > 120 {
			return nil, fmt.Errorf("validity in months must be between 1 and 120")
		}
		if *p.ValidityInMonths != 0 {
			doc.ValidityInMonths = *p.ValidityInMonths
		}
	}

	doc.LifetimeTrigger.setDefault()
	if err := doc.LifetimeTrigger.fromInput(p.LifetimeTrigger, doc.ValidityInMonths); err != nil {
		return nil, err
	}

	return doc, nil
}
