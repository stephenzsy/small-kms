package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
)

func getRootCaRefs() []RefWithMetadata {
	return []RefWithMetadata{
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_RootCA, DisplayName: "Root CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestRootCA, DisplayName: "Test Root CA", Type: RefTypeNamespace},
	}
}

func getIntCaRefs() []RefWithMetadata {
	return []RefWithMetadata{
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAService, DisplayName: "Services Intermediate CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAIntranet, DisplayName: "Intranet Intermediate CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAAadSp, DisplayName: "AAD Client Secret Intermediate CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestIntCA, DisplayName: "Test Intermediate CA", Type: RefTypeNamespace},
	}
}

func (s *adminServer) ListNamespacesByTypeV2(c *gin.Context, nsType NamespaceTypeShortName) {
	if !authAdminOnly(c) {
		return
	}
	var odataType graph.MsGraphOdataType
	switch nsType {
	case NSTypeRootCA:
		c.JSON(http.StatusOK, getRootCaRefs())
		return
	case NSTypeIntCA:
		c.JSON(http.StatusOK, getIntCaRefs())
		return
	}
	switch nsType {
	case NSTypeGroup:
		odataType = graph.MsGraphOdataTypeGroup
	case NSTypeUser:
		odataType = graph.MsGraphOdataTypeUser
	case NSTypeServicePrincipal:
		odataType = graph.MsGraphOdataTypeServicePrincipal
	case NSTypeDevice:
		odataType = graph.MsGraphOdataTypeDevice
	case NSTypeApplication:
		odataType = graph.MsGraphOdataTypeApplication
	default:
		respondPublicErrorMsg(c, http.StatusBadRequest, "unsupported namespace type")
		return
	}
	dirObjs, err := s.graphService.ListGraphProfilesByType(c, odataType)

	if err != nil {
		common.RespondError(c, err)
		return
	}

	r := make([]RefWithMetadata, len(dirObjs))
	for i, doc := range dirObjs {
		profileDocPopulateRefWithMetadata(&doc, &r[i])
		r[i].Type = RefTypeNamespace
	}

	c.JSON(http.StatusOK, r)
}
