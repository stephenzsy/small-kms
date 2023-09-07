package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) ApplyPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}

	// load policy
	policy, err := s.GetPolicyDoc(c, namespaceID, policyID)
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, nil)
		} else {
			log.Error().Err(err).Msg("Internal error")
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	b := ApplyPolicyRequest{}
	err = c.BindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	switch policy.PolicyType {
	case PolicyTypeCertRequest:
		section := policy.CertRequest
		if section == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid policy to request certificate"})
			return
		}
		shouldRenew, stateDoc, reason, err := section.evaluateForAction(c, s, namespaceID, policy, b.ForceRenewCertificate)
		if err != nil {
			log.Error().Err(err).Msg("Internal error")
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}
		log.Info().Msgf("shouldRenew: %v, reason: %s", shouldRenew, reason)
		if shouldRenew {
			stateDoc, err = section.action(c, s, namespaceID, policy)
			if err != nil {
				log.Error().Err(err).Msg("Internal error")
				c.JSON(500, gin.H{"error": "internal error"})
				return
			}
		}
		c.JSON(200, stateDoc.ToPolicyState())
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
}
