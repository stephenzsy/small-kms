package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) SyncNamespaceInfoV2(c *gin.Context, namespaceId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	graphClient, err := s.msGraphClient(c)
	if err != nil {
		common.RespondError(c, err)
		return
	}

	obj, err := graphClient.DirectoryObjects().ByDirectoryObjectId(namespaceId.String()).Get(c, nil)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("directoryObject:%s", namespaceId))
		if errors.Is(err, common.ErrStatusNotFound) {
			// delete the doc if it exists
			profileDoc, err := s.graphService.GetGraphProfileDoc(c, namespaceId, graph.MsGraphOdataTypeAny)
			if err == nil {
				if err = kmsdoc.AzCosmosDelete(c, s.AzCosmosContainerClient(), profileDoc); err != nil {
					respondInternalError(c, err, fmt.Sprintf("failed to delete directory object in cosmos: %s", namespaceId))
					return
				}
			}
			respondPublicError(c, http.StatusNotFound, err)
			return
		}
		respondInternalError(c, err, fmt.Sprintf("failed to get directory object: %s", namespaceId))
		return
	}

	profilableObj, ok := obj.(graph.GraphProfileable)
	if !ok {
		respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("directory object is not supported: %s", namespaceId))
		return
	}
	profileDoc := s.graphService.NewGraphProfileDoc(s.TenantID(), profilableObj)
	if profileDoc.IsValid() {
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), profileDoc); err != nil {
			respondInternalError(c, err, fmt.Sprintf("failed to upsert directory object in cosmos: %s", namespaceId))
			return
		}
	} else {
		respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("directory object is not supported: %s", namespaceId))
		return
	}

	c.JSON(http.StatusOK, newNamespaceInfoFromProfileDoc(profileDoc))
}

func newNamespaceInfoFromProfileDoc(doc graph.GraphProfileDocument) *NamespaceInfo {
	if doc == nil {
		return nil
	}
	p := new(NamespaceInfo)
	profileDocPopulateRefWithMetadata(doc, &p.Ref)
	p.Ref.Type = RefTypeNamespace
	p.ObjectType = OdataTypeToNSType(doc.GetOdataType())
	return p
}
