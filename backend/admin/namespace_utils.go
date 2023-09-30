package admin

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

// Deprecated
func IsNamespaceManagementAdminRequired(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case common.WellKnownID_RootCA,
		common.WellKnownID_IntCAService,
		common.WellKnownID_IntCAIntranet,
		common.WellKnownID_IntCAAadSp,
		common.WellKnownID_TestRootCA:
		return true
	}
	return false
}

var rootCaAllowedIDs = map[uuid.UUID]bool{
	common.WellKnownID_RootCA:     true,
	common.WellKnownID_TestRootCA: true,
}

func isAllowedRootCaNamespace(namespaceID uuid.UUID) bool {
	return rootCaAllowedIDs[namespaceID]
}

var intermediateCaAllowedIDs = map[uuid.UUID]bool{
	common.WellKnownID_IntCAService:  true,
	common.WellKnownID_IntCAIntranet: true,
	common.WellKnownID_IntCAAadSp:    true,
	common.WellKnownID_TestIntCA:     true,
}

func isAllowedIntCaNamespace(namespaceID uuid.UUID) bool {
	return intermediateCaAllowedIDs[namespaceID]
}

func isAllowedCaNamespace(namespaceID uuid.UUID) bool {
	return isAllowedRootCaNamespace(namespaceID) || isAllowedIntCaNamespace(namespaceID)
}

var testCaIDs = map[uuid.UUID]bool{
	common.WellKnownID_TestRootCA: true,
	common.WellKnownID_TestIntCA:  true,
}

func isTestCA(namespaceID uuid.UUID) bool {
	return testCaIDs[namespaceID]
}

// returns a tuple of (isValid, needs graph validation)
func validateNamespaceType(nsType NamespaceTypeShortName, nsID uuid.UUID) (bool, bool) {
	switch nsType {
	case NSTypeRootCA:
		return isAllowedRootCaNamespace(nsID), false
	case NSTypeIntCA:
		return isAllowedIntCaNamespace(nsID), false
	case NSTypeServicePrincipal,
		NSTypeGroup,
		NSTypeDevice,
		NSTypeUser,
		NSTypeApplication:
		return nsID.Version() == 4, true
	}
	return false, false
}

func validateNamespaceTypeWithDirDoc(nsType NamespaceTypeShortName, doc *DirectoryObjectDoc) bool {
	if doc == nil {
		return false
	}
	switch nsType {
	case NSTypeServicePrincipal:
		return doc.OdataType == string(NamespaceTypeMsGraphServicePrincipal)
	case NSTypeGroup:
		return doc.OdataType == string(NamespaceTypeMsGraphGroup)
	case NSTypeDevice:
		return doc.OdataType == string(NamespaceTypeMsGraphDevice)
	case NSTypeUser:
		return doc.OdataType == string(NamespaceTypeMsGraphUser)
	case NSTypeApplication:
		return doc.OdataType == string(NamespaceTypeMsGraphApplication)
	}
	return false
}
