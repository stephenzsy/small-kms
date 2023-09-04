package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) ApplyPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	_, ok := authNamespaceAdminOrSelf(c, namespaceID)
	if !ok {
		return
	}

	// load policy
	policy, err := s.GetPolicyDoc(c, namespaceID, policyID)
	if err != nil {
		if kmsdoc.IsNotFound(err) {
			c.JSON(404, gin.H{"error": "not found"})
		} else {
			log.Error().Err(err).Msg("Internal error")
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	b := ApplyPolicyRequest{}
	err = c.BindJSON(&b)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid apply policy parameter"})
		return
	}
	switch policy.PolicyType {
	case PolicyTypeCertRequest:
		section := policy.CertRequest
		if section == nil {
			c.JSON(400, gin.H{"error": "invalid policy to request certificate"})
			return
		}
		shouldRenew, reason, err := section.evaluateForAction(c, s, namespaceID, policy, b.ForceRenewCertificate)
		if err != nil {
			log.Error().Err(err).Msg("Internal error")
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}
		log.Info().Msgf("shouldRenew: %v, reason: %s", shouldRenew, reason)
		if shouldRenew {
			resultDoc, err := section.action(c, s, namespaceID, policy)
			if err != nil {
				log.Error().Err(err).Msg("Internal error")
				c.JSON(500, gin.H{"error": "internal error"})
				return
			}
			c.JSON(200, resultDoc.ToPolicyState())
		} else {
			c.JSON(200, gin.H{"message": reason})
		}
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
}
