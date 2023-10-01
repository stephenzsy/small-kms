package admin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/deviceswithdeviceid"
)

func (s *adminServer) GetMyProfilesV1(c *gin.Context) {
	/*
		r := make([]*NamespaceProfile, 0, 2)

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
		if len(callerDeviceId) > 0 {
			// resolve
			deviceOid, err := s.resolveObjectIDFromDeviceID(c, callerDeviceId)
			if err != nil {
				// failed to resove device id, ignore
				log.Err(err).Msgf("Failed to resolve device id: %s", callerDeviceId)
			} else if deviceOid != callerPrincipalId {
				profile, err := s.GetNamespaceProfile(c, deviceOid)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get namespace profile for device: %s", callerDeviceId)
					c.JSON(500, gin.H{"message": "internal error"})
					return
				}
				if profile != nil {
					r = append(r, profile)
				}
			}
		}*/

	c.JSON(404, nil)
}

func (s *adminServer) resolveObjectIDFromDeviceID(c context.Context, deviceID string) (uuid.UUID, error) {
	device, err := s.msGraphClient.DevicesWithDeviceId(&deviceID).Get(c,
		&deviceswithdeviceid.DevicesWithDeviceIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceswithdeviceid.DevicesWithDeviceIdRequestBuilderGetQueryParameters{
				Select: []string{"id"},
			},
		})
	if err != nil {
		return uuid.Nil, err
	}
	idStr := *device.GetId()
	return uuid.Parse(idStr)
}

func (s *adminServer) SyncMyProfilesV1(c *gin.Context) {
	/*
		r := make([]*NamespaceProfile, 0, 2)

		callerPrincipalId := auth.CallerPrincipalId(c)
		if callerPrincipalId != uuid.Nil {
			profile, status, err := s.RegisterNamespaceProfile(c, callerPrincipalId)
			if err != nil {
				if status == http.StatusInternalServerError {
					log.Error().Err(err).Msg("Failed to register graph object")
					c.JSON(500, gin.H{"message": "internal error"})
					return
				}
			}
			if profile != nil && status == http.StatusOK {
				r = append(r, profile)
			}
		}
		callerDeviceId := auth.CallerPrincipalDeviceID(c)
		if len(callerDeviceId) > 0 {
			deviceOid, err := s.resolveObjectIDFromDeviceID(c, callerDeviceId)
			if err != nil {
				// failed to resove device id, ignore
				log.Err(err).Msgf("Failed to resolve device id: %s", callerDeviceId)
			} else if deviceOid != callerPrincipalId {
				profile, status, err := s.RegisterNamespaceProfile(c, deviceOid)
				if err != nil {
					if status == http.StatusInternalServerError {
						log.Error().Err(err).Msg("Failed to register graph object")
						c.JSON(500, gin.H{"message": "internal error"})
						return
					}
				}
				if profile != nil && status == http.StatusOK {
					r = append(r, profile)
				}
			}
		}

		c.JSON(200, r)
	*/
	c.JSON(404, nil)
}
