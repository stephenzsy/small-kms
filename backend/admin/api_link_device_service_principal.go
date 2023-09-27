package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) LinkDeviceServicePrincipalV2(c *gin.Context, namespaceId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	dirDoc, err := s.syncDirDoc(c, namespaceId)
	if err != nil {
		if common.IsGraphODataErrorNotFound(err) || common.IsAzNotFound(err) {
			respondPublicError(c, http.StatusNotFound, err)
			return
		}
		respondInternalError(c, err, "failed to sync directory")
		return
	}

	if dirDoc.OdataType != string(NamespaceTypeMsGraphDevice) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "namespace is not a device"})
		return
	}

	if dirDoc.Device != nil && dirDoc.Device.LinkedApplicationClientID != nil {
		c.JSON(http.StatusAccepted, DeviceServicePrincipal{
			ApplicationClientID:      *dirDoc.Device.LinkedApplicationClientID,
			ApplicationObjectID:      *dirDoc.Device.LinkedApplicationObjectID,
			ServicePrincipalObjectID: *dirDoc.Device.LinkedServicePrincipalObjectID,
			DeviceId:                 dirDoc.Device.DeviceID,
		})
		return
	}
}
