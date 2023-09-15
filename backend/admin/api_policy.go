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

func getBuiltinPolicyRefs(namespaceID uuid.UUID) []PolicyRef {
	switch {
	case IsRootCANamespace(namespaceID):
		return []PolicyRef{
			{NamespaceID: namespaceID, ID: namespaceID, PolicyType: PolicyTypeCertRequest},
		}
	case IsIntCANamespace(namespaceID):
		rootCaNs := wellKnownNamespaceID_RootCA
		if IsTestCA(namespaceID) {
			rootCaNs = testNamespaceID_RootCA
		}
		return []PolicyRef{
			{NamespaceID: namespaceID, ID: rootCaNs, PolicyType: PolicyTypeCertRequest},
		}

	}
	return nil
}

func (s *adminServer) ListPoliciesV1(c *gin.Context, namespaceID uuid.UUID) {
	// validate
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}
	builtInList := getBuiltinPolicyRefs(namespaceID)
	if builtInList != nil {
		c.JSON(http.StatusOK, builtInList)
		return
	}
	l, err := s.ListPoliciesByNamespace(c, namespaceID)
	if err != nil {
		log.Err(err).Msg("Internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	results := make([]PolicyRef, len(l))
	for i, item := range l {
		item.PopulatePolicyRef(&results[i])
	}
	c.JSON(http.StatusOK, results)
}

func (s *adminServer) PutPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	// validate
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	p := PolicyParameters{}

	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	policyDoc := new(PolicyDoc)
	var dirProfile *DirectoryObjectDoc
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
			dirProfile, err := s.GetDirectoryObjectDoc(c, namespaceID)
			if err != nil {
				if common.IsAzNotFound(err) {
					c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("namespace not registered yet: %s", namespaceID)})
					return
				}
				log.Error().Err(err).Msg("failed to get directory profile")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
			switch dirProfile.OdataType {
			case string(NamespaceTypeMsGraphGroup),
				string(NamespaceTypeMsGraphServicePrincipal):
				// ok
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "namespace not supported yet"})
				return
			}
		}
		docSection := new(PolicyCertRequestDocSection)
		if err := docSection.validateAndFillWithParameters(p.CertRequest, namespaceID, dirProfile); err != nil {
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
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
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
