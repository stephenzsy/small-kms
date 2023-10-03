package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
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

func validateTemplateIdentifiers(nsType NamespaceTypeShortName, nsID uuid.UUID, templateID uuid.UUID, name string) (string, bool) {
	if templateID.Version() == 4 {
		// only allow non default prefixed name for user specified template ID
		return name, (name != "" && !strings.HasPrefix(name, "default"))
	} else {
		// assign default template name for well known IDs
		switch {
		case templateID == uuid.Nil:
			// allow default template for specified namespace type
			switch nsType {
			case NSTypeRootCA,
				NSTypeIntCA,
				NSTypeServicePrincipal:
				return string(common.DefaultCertTemplateName_GlobalDefault), true
			}
		case nsType == NSTypeGroup &&
			templateID == common.GetCanonicalCertificateTemplateID(nsID, common.DefaultCertTemplateName_ServicePrincipalClientCredential):
			return string(common.DefaultCertTemplateName_ServicePrincipalClientCredential), true
		}
	}
	return "invalid", false
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
	if fixedName, ok := validateTemplateIdentifiers(nsType, nsID, templateId, displayName); ok {
		displayName = fixedName
	} else {
		return nil, fmt.Errorf("template ID %s is not valid for namespace type %s", templateId, nsType)
	}

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
		if p.Usage != UsageRootCA {
			return nil, errors.New("root ca must be used for root ca certificate")
		}
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
		if p.Usage != UsageIntCA {
			return nil, errors.New("intermediate ca must be used for intermediate ca certificate")
		}
		doc.Usage = UsageIntCA
		doc.ValidityInMonths = 36 // default 3 years

	default:
		if !isAllowedIntCaNamespace(p.Issuer.NamespaceID) {
			return nil, fmt.Errorf("service principal/group issuer namespace ID %s is not an intermediate ca namespace ID", p.Issuer.NamespaceID)
		}
		p.populateDocIssuer(doc, NSTypeIntCA)
		if p.Usage == UsageRootCA || p.Usage == UsageIntCA {
			return nil, errors.New("root ca/intermediate ca must not be used for other certificate")
		}
		if p.Usage == UsageAADClientCredential && nsType != NSTypeServicePrincipal {
			return nil, errors.New("AAD client credential certificate must be used for service principal")
		}
		doc.Usage = p.Usage
		doc.ValidityInMonths = 12 // default 1 year
	}

	// validate and populate key properties
	doc.KeyProperties.setDefault()
	doc.KeyProperties.ReuseKey = p.ReuseKey
	switch nsType {
	case NSTypeRootCA,
		NSTypeIntCA:
		if isTestCA(nsID) {
			doc.KeyProperties.setECDSA(CurveNameP384)
		} else {
			doc.KeyProperties.setRSA(AlgRS384, KeySize4096)
		}
		// ignore input key properties for CA
	default:
		if err := doc.KeyProperties.fromJwkProperties(p.KeyProperties); err != nil {
			return nil, err
		}
		if p.Usage == UsageAADClientCredential && doc.KeyProperties.Kty == KeyTypeEC {
			return nil, fmt.Errorf("AAD client credential certificate must use RSA key")
		}
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

	doc.Subject = CertificateTemplateDocSubject{CertificateSubject: p.Subject}
	doc.SubjectAlternativeNames = sanitizeSANs(p.SubjectAlternativeNames)

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
