package common

import (
	"github.com/google/uuid"
)

type WellKnownID uuid.UUID

type WellKnownIdentifier int

const (
	IdentifierUnknown WellKnownIdentifier = iota

	IdentifierRootCA
	IdentifierTestRootCA

	IdentifierIntCAService
	IdentifierIntCAIntranet
	IdentifierTestIntCA

	IdentifierDirectory
	IdentifierIssuerPolicy
)

var (
	// nil is reserved

	// root CA --1 ~ --f
	idRootCA     = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	idTestRootCA = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-00000000000f"))

	// intermediate CAs --10 ~ --ff
	idIntCAService  = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-000000000011"))
	idIntCAIntranet = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-000000000012"))
	idTestIntCA     = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-0000000000ff"))

	idDirectory = WellKnownID(uuid.MustParse(MustGetenv("AZURE_TENANT_ID")))

	// special policy ids --1-1 ~ --1-f
	idIssuerPolicy = WellKnownID(uuid.MustParse("00000000-0000-0000-0001-000000000001"))
)

var idMap = map[WellKnownIdentifier]WellKnownID{
	IdentifierRootCA:     idRootCA,
	IdentifierTestRootCA: idTestRootCA,

	IdentifierIntCAService:  idIntCAService,
	IdentifierIntCAIntranet: idIntCAIntranet,
	IdentifierTestIntCA:     idTestIntCA,

	IdentifierDirectory:    idDirectory,
	IdentifierIssuerPolicy: idIssuerPolicy,
}

func GetID(identifier WellKnownIdentifier) uuid.UUID {
	return uuid.UUID(idMap[identifier])
}
