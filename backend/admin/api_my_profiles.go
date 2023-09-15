package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
)

func (s *adminServer) GetMyProfilesV1(c *gin.Context) {
	r := make([]*NamespaceProfile, 2)

	callerPrincipalId := auth.CallerPrincipalId(c)
	if callerPrincipalId != uuid.Nil {
		profile, err := s.GetNamespaceProfile(c, callerPrincipalId)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get namespace profile for caller %s", callerPrincipalId.String())
			c.JSON(500, gin.H{"message": "internal error"})
			return
		}
		if profile != nil {
			r = append(r, profile)
		}
	}
	callerDeviceId := auth.CallerPrincipalDeviceID(c)
	if callerDeviceId != uuid.Nil && callerDeviceId != callerPrincipalId {
		profile, err := s.GetNamespaceProfile(c, callerDeviceId)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get namespace profile for device: %s", callerDeviceId.String())
			c.JSON(500, gin.H{"message": "internal error"})
			return
		}
		if profile != nil {
			r = append(r, profile)
		}
	}

	c.JSON(200, r)
}
