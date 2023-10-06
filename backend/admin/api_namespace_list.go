package admin

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

var (
	rootCaRef     = RefWithMetadata{NamespaceID: uuid.Nil, ID: common.WellKnownID_RootCA, DisplayName: "Root CA", Type: RefTypeNamespace}
	rootTestCaRef = RefWithMetadata{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestRootCA, DisplayName: "Test Root CA", Type: RefTypeNamespace}
)

func getRootCaRefs() []RefWithMetadata {
	return []RefWithMetadata{
		rootCaRef,
		rootTestCaRef,
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
