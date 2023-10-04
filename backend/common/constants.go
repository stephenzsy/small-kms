package common

import (
	"github.com/google/uuid"
)

type WellKnownID uuid.UUID

type WellKnownIdentifier int

var (
	// nil is reserved

	// root CA --1 ~ --f
	WellKnownID_RootCAMin  = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	WellKnownID_RootCA     = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	WellKnownID_TestRootCA = uuid.MustParse("00000000-0000-0000-0000-00000000000f")
	WellKnownID_RootCAMax  = uuid.MustParse("00000000-0000-0000-0000-00000000000f")

	// intermediate CAs --10 ~ --ff
	WellKnownID_IntCAService  = uuid.MustParse("00000000-0000-0000-0000-000000000011")
	WellKnownID_IntCAIntranet = uuid.MustParse("00000000-0000-0000-0000-000000000012")
	WellKnownID_IntCAAadSp    = uuid.MustParse("00000000-0000-0000-0000-000000000013")
	WellKnownID_TestIntCA     = uuid.MustParse("00000000-0000-0000-0000-0000000000ff")
)

type WellKnownCertTemplateName string

const (
	DefaultCertTemplateName_GlobalDefault                    WellKnownCertTemplateName = "default"
	DefaultCertTemplateName_ServicePrincipalClientCredential WellKnownCertTemplateName = "default-service-principal-client-credential"
)
const (
	NSRelNameDASPLink = "device-application-service-principal-link"
)
