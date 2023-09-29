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

func IsRootCANamespace(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case common.WellKnownID_RootCA,
		common.WellKnownID_TestRootCA:
		return true
	}
	return false
}

func IsIntCANamespace(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case common.WellKnownID_IntCAService,
		common.WellKnownID_IntCAIntranet,
		common.WellKnownID_IntCAAadSp,
		common.WellKnownID_TestIntCA:
		return true
	}
	return false
}

func IsCANamespace(namespaceID uuid.UUID) bool {
	return IsRootCANamespace(namespaceID) || IsIntCANamespace(namespaceID)
}

func IsTestCA(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case common.WellKnownID_TestRootCA,
		common.WellKnownID_TestIntCA:
		return true
	}
	return false
}

// returns a tuple of (isValid, needs graph validation)
func validateNamespaceType(nsType NamespaceTypeShortName, nsID uuid.UUID) (bool, bool) {
	switch nsType {
	case NSTypeRootCA:
		return IsRootCANamespace(nsID), false
	case NSTypeIntCA:
		return IsIntCANamespace(nsID), false
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
