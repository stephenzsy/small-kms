package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) SyncNamespaceInfoV2(c *gin.Context, namespaceType NamespaceTypeShortName, namespaceId uuid.UUID) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceId); !ok {
		return
	}
	isValid, isGraphValidationNeeded := validateNamespaceType(namespaceType, namespaceId)
	if !isValid {
		respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace type %s is not valid for ID: %s", namespaceType, namespaceId))
		return
	}
	if !isGraphValidationNeeded {
		respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace type %s is invalid for sync", namespaceType))
		return
	}
	doc, err := s.genDirDocFromMsGraph(c, namespaceId)
	if err != nil {
		if common.IsGraphODataErrorNotFound(err) || common.IsAzNotFound(err) {
			respondPublicError(c, http.StatusNotFound, err)
			return
		}
		respondInternalError(c, err, fmt.Sprintf("failed to sync directory object: %s", namespaceId))
		return
	}
	if validateNamespaceTypeWithDirDoc(namespaceType, doc) {
		err := kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, doc)
		if err != nil {
			respondInternalError(c, err, fmt.Sprintf("failed to upsert directory object in cosmos: %s", namespaceId))
			return
		}
	} else {
		respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace type %s is not valid for ID: %s", namespaceType, namespaceId))
		return
	}

	c.JSON(http.StatusOK, doc.toNamespaceInfo())
}
