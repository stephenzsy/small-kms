package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
)

func (s *adminServer) GetNamespaceInfoV2(c *gin.Context, namespaceId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	if isAllowedRootCaNamespace(namespaceId) {

		for _, r := range getRootCaRefs() {
			if r.ID == namespaceId {
				c.JSON(http.StatusOK, &NamespaceInfo{
					Ref:        r,
					ObjectType: NSTypeRootCA,
				})
				return
			}
		}
	} else if isAllowedIntCaNamespace(namespaceId) {
		for _, r := range getIntCaRefs() {
			if r.ID == namespaceId {
				c.JSON(http.StatusOK, &NamespaceInfo{
					Ref:        r,
					ObjectType: NSTypeIntCA,
				})
				return
			}
		}
	}

	profileDoc, err := s.graphService.GetGraphProfileDoc(c, namespaceId, graph.MsGraphOdataTypeAny)
	if err != nil {
		common.RespondError(c, err)
	}

	c.JSON(http.StatusOK, newNamespaceInfoFromProfileDoc(profileDoc))
}
