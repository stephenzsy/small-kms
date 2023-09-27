package admin

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

var wellKnownNamespaceID_IntCAService = common.GetID(common.IdentifierIntCAService)
var wellKnownNamespaceID_IntCaIntranet uuid.UUID = common.GetID(common.IdentifierIntCAIntranet)

var wellknownNamespaceID_directoryID = common.GetID(common.IdentifierDirectory)

var testNamespaceID_IntCA = common.GetID(common.IdentifierTestIntCA)

func IsNamespaceManagementAdminRequired(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case common.WellKnownID_RootCA,
		wellKnownNamespaceID_IntCAService,
		wellKnownNamespaceID_IntCaIntranet,
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
	case wellKnownNamespaceID_IntCAService,
		wellKnownNamespaceID_IntCaIntranet,
		testNamespaceID_IntCA:
		return true
	}
	return false
}

func IsCANamespace(namespaceID uuid.UUID) bool {
	return IsRootCANamespace(namespaceID) || IsIntCANamespace(namespaceID)
}

func IsTestCA(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case common.WellKnownID_TestRootCA, testNamespaceID_IntCA:
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
