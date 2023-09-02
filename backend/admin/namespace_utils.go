package admin

import "github.com/google/uuid"

var wellKnownNamespaceID_RootCA uuid.UUID = uuid.MustParse(string(WellKnownNamespaceIDStrRootCA))
var wellKnownNamespaceID_IntCAService uuid.UUID = uuid.MustParse(string(WellKnownNamespaceIDStrIntCAService))
var wellKnownNamespaceID_IntCaIntranet uuid.UUID = uuid.MustParse(string(WellKnownNamespaceIDStrIntCAIntranet))

var namespacePrefixMapping = map[uuid.UUID]string{
	wellKnownNamespaceID_RootCA:        "root-ca-",
	wellKnownNamespaceID_IntCAService:  "int-ca-service-",
	wellKnownNamespaceID_IntCaIntranet: "int-ca-intranet-",
	testNamespaceID_RootCA:             "test-root-ca-",
}

var testNamespaceID_RootCA uuid.UUID = uuid.MustParse(string(TestNamespaceIDStrRootCA))

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
		wellKnownNamespaceID_IntCaIntranet:
		return true
	}
	return false
}

func IsCANamespace(namespaceID uuid.UUID) bool {
	return IsRootCANamespace(namespaceID) || IsIntCANamespace(namespaceID)
}
