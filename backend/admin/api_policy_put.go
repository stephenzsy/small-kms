package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func isPolicyTypeValidForId(policyType PolicyType, policyID uuid.UUID) bool {
	switch policyID {
	case defaultPolicyIdCertRequest:
		return policyType == PolicyTypeCertRequest
	case defaultPolicyIdCertEnroll:
		return policyType == PolicyTypeCertEnroll
	case defaultPolicyIdCertAadAppCredential:
		return policyType == PolicyTypeCertAadAppClientCredential
	}
	return policyID.Version() == 4
}

func (s *adminServer) PutPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyIdentifier string) {
	// validate
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	policyID, err := resolvePolicyIdentifier(policyIdentifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid policy identifier: %s", policyIdentifier)})
		return
	}

	p := PolicyParameters{}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if !isPolicyTypeValidForId(p.PolicyType, policyID) {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("policy type %s is not valid for policy id %s", p.PolicyType, policyID)})
		return
	}

	policyDoc := new(PolicyDoc)
	switch p.PolicyType {
	case PolicyTypeCertEnroll:
		// verify namespace is an intermediate CA
		if !IsIntCANamespace(namespaceID) {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Issuer namespace is invalid for certificate enrollment: %s", namespaceID.String())})
			return
		}
		docSection := new(PolicyCertEnrollDocSection)
		if len(p.CertEnroll.AllowedUsages) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "allowed usages must be specified"})
			return
		}
		docSection.AllowedUsages = p.CertEnroll.AllowedUsages
		if p.CertEnroll.MaxValidityInMonths < 1 || p.CertEnroll.MaxValidityInMonths > 120 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "max validity in months must be between 1 and 120, inclusive"})
			return
		}
		docSection.MaxValidityInMonths = p.CertEnroll.MaxValidityInMonths
		policyDoc.CertEnroll = docSection
	case PolicyTypeCertAadAppClientCredential:
		dirDoc, err := s.getDirectoryObjectDoc(c, namespaceID)
		if err != nil {
			if common.IsAzNotFound(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("namespace not registered yet: %s", namespaceID)})
				return
			}
			log.Error().Err(err).Msg("failed to get directory profile")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if dirDoc.OdataType != string(NamespaceTypeMsGraphApplication) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("policy can only be registered on application: %s", namespaceID)})
			return
		}
		docSection := new(PolicyCertAadAppCredDocSection)
		if err := docSection.validateAndFillWithParameters(p.CertAadAppCred); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		policyDoc.CertAadAppCred = docSection
	default:
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	policyDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicy, policyID)
	policyDoc.NamespaceID = namespaceID
	policyDoc.PolicyType = p.PolicyType

	// write to DB
	if err := kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, policyDoc); err != nil {
		log.Printf("Internal error: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, policyDoc.ToPolicy())
}
