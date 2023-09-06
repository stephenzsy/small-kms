package admin

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

var wellKnownNamespaceID_RootCA = common.GetID(common.IdentifierRootCA)
var wellKnownNamespaceID_IntCAService = common.GetID(common.IdentifierIntCAService)
var wellKnownNamespaceID_IntCaIntranet uuid.UUID = common.GetID(common.IdentifierIntCAIntranet)

var directoryID = common.GetID(common.IdentifierDirectory)

var testNamespaceID_RootCA = common.GetID(common.IdentifierTestRootCA)
var testNamespaceID_IntCA = common.GetID(common.IdentifierTestIntCA)

func IsNamespaceManagementAdminRequired(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case wellKnownNamespaceID_RootCA,
		wellKnownNamespaceID_IntCAService,
		wellKnownNamespaceID_IntCaIntranet,
		testNamespaceID_RootCA:
		return true
	}
	return false
}

func IsRootCANamespace(namespaceID uuid.UUID) bool {
	switch namespaceID {
	case wellKnownNamespaceID_RootCA,
		testNamespaceID_RootCA:
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
	case testNamespaceID_RootCA, testNamespaceID_IntCA:
		return true
	}
	return false
}
