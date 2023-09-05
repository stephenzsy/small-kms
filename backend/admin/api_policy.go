package admin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) PutPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	// validate
	_, ok := authNamespaceAdminOrSelf(c, namespaceID)
	if !ok {
		return
	}

	p := PolicyParameters{}

	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	policyDoc := new(PolicyDoc)
	switch p.PolicyType {
	case PolicyTypeCertRequest:
		switch {
		case IsRootCANamespace(namespaceID):
			if namespaceID != policyID {
				c.JSON(http.StatusForbidden, gin.H{"message": "root namespace must have policy name as the same as the namespace id"})
				return
			}
		case IsIntCANamespace(namespaceID):
			if IsTestCA(namespaceID) {
				if policyID != testNamespaceID_RootCA {
					c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Issuer %s does not allow the requester namespace: %s", policyID.String(), namespaceID.String())})
					return
				}
			} else {
				if policyID != wellKnownNamespaceID_RootCA {
					c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Issuer %s does not allow the requester namespace: %s", policyID.String(), namespaceID.String())})
					return
				}
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "namespace not supported yet"})
			return
		}
		docSection := new(PolicyCertRequestDocSection)
		if err := docSection.validateAndFillWithParameters(p.CertRequest, namespaceID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		policyDoc.CertRequest = docSection
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

func (s *adminServer) GetPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	// validate
	_, ok := authNamespaceAdminOrSelf(c, namespaceID)
	if !ok {
		return
	}
	pd, err := s.GetPolicyDoc(c, namespaceID, policyID)
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, nil)
		} else {
			log.Printf("Internal error: %s", err.Error())
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(200, pd.ToPolicy())
}
