package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func getRootCaRefs() []RefWithMetadata {
	return []RefWithMetadata{
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_RootCA, Metadata: map[string]string{RefPropertyKeyDisplayName: "Root CA"}, Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestRootCA, Metadata: map[string]string{RefPropertyKeyDisplayName: "Test Root CA"}, Type: RefTypeNamespace},
	}
}

func getIntCaRefs() []RefWithMetadata {
	return []RefWithMetadata{
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAService, Metadata: map[string]string{RefPropertyKeyDisplayName: "Services Intermediate CA"}, Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAIntranet, Metadata: map[string]string{RefPropertyKeyDisplayName: "Intranet Intermediate CA"}, Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAAadSp, Metadata: map[string]string{RefPropertyKeyDisplayName: "AAD Client Secret Intermediate CA"}, Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestIntCA, Metadata: map[string]string{RefPropertyKeyDisplayName: "Test Intermediate CA"}, Type: RefTypeNamespace},
	}
}

func (s *adminServer) ListNamespacesByTypeV2(c *gin.Context, nsType NamespaceTypeShortName) {
	if !authAdminOnly(c) {
		return
	}
	switch nsType {
	case NSTypeRootCA:
		c.JSON(http.StatusOK, getRootCaRefs())
		return
	case NSTypeIntCA:
		c.JSON(http.StatusOK, getIntCaRefs())
		return
	}
	var odType string
	switch nsType {
	case NSTypeGroup:
		odType = "#microsoft.graph.group"
	case NSTypeUser:
		odType = "#microsoft.graph.user"
	case NSTypeServicePrincipal:
		odType = "#microsoft.graph.servicePrincipal"
	case NSTypeDevice:
		odType = "#microsoft.graph.device"
	case NSTypeApplication:
		odType = "#microsoft.graph.application"
	default:
		respondPublicErrorMsg(c, http.StatusBadRequest, "unsupported namespace type")
		return
	}
	dirObjs, err := s.listDirectoryObjectByType(c, odType)

	if err != nil {
		respondInternalError(c, err, "failed to list directory objects")
		return
	}
	r := make([]RefWithMetadata, len(dirObjs))
	for i, doc := range dirObjs {
		baseDocPopulateRefWithMetadata(&doc.BaseDoc, &r[i], nsType)
		r[i].Metadata = map[string]string{RefPropertyKeyDisplayName: doc.DisplayName}
		r[i].Type = RefTypeNamespace
	}

	c.JSON(http.StatusOK, r)
}
