package admin

import "github.com/google/uuid"

var wellKnownNamespaceID_RootCA uuid.UUID = uuid.MustParse(string(WellKnownNamespaceIDStrRootCA))
var wellKnownNamespaceID_IntCAService uuid.UUID = uuid.MustParse(string(WellKnownNamespaceIDStrIntCAService))
var wellKnownNamespaceID_IntCAClient uuid.UUID = uuid.MustParse(string(WellKnownNamespaceIDStrIntCAClient))

var namespacePrefixMapping = map[uuid.UUID]string{
	wellKnownNamespaceID_RootCA:       "root-ca-",
	wellKnownNamespaceID_IntCAService: "int-ca-service-",
	wellKnownNamespaceID_IntCAClient:  "int-ca-client-",
}
