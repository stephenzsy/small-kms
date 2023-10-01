package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) ListPoliciesV1(c *gin.Context, namespaceID uuid.UUID) {
	// validate
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}
	l, err := s.listPoliciesByNamespace(c, namespaceID)
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

var (
	defaultPolicyIdCertEnroll = common.GetID(common.DefaultPolicyIdCertEnroll)
)

func resolvePolicyIdentifier(policyIdentifier string) (uuid.UUID, error) {
	switch policyIdentifier {
	case string(PolicyTypeCertEnroll):
		return defaultPolicyIdCertEnroll, nil
	}
	return uuid.Parse(policyIdentifier)
}

func (s *adminServer) GetPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyIdentifier string) {
	// validate
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}
	policyID, err := resolvePolicyIdentifier(policyIdentifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid policy identifier: %s", policyIdentifier)})
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

/*

// Delete Certificate Policy
// (DELETE /v1/{namespaceId}/policies/{policyIdentifier})
func (s *adminServer) DeletePolicyV1(c *gin.Context, namespaceID uuid.UUID, policyIdentifier string, params DeletePolicyV1Params) {
	purge := false
	if params.Purge != nil && *params.Purge {
		if !auth.CallerPrincipalHasAdminRole(c) {
			c.JSON(http.StatusForbidden, gin.H{"message": "only admin can purge"})
			return
		}
		purge = true
	}
	// validate
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}
	policyID, err := resolvePolicyIdentifier(policyIdentifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid policy identifier: %s", policyIdentifier)})
		return
	}
	err = s.deletePolicyDoc(c, namespaceID, policyID, purge)
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, nil)
		} else {
			log.Printf("Internal error: %s", err.Error())
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	if purge {
		c.JSON(204, nil)
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
*/
