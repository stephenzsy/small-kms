package admin

import (
	"github.com/google/uuid"
)

func isPolicyTypeValidForId(policyType PolicyType, policyID uuid.UUID) bool {
	switch policyID {
	case defaultPolicyIdCertEnroll:
		return policyType == PolicyTypeCertEnroll
	}
	return policyID.Version() == 4
}

/*
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
		if !isAllowedIntCaNamespace(namespaceID) {
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
*/
