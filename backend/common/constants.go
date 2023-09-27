package common

import (
	"github.com/google/uuid"
)

type WellKnownID uuid.UUID

type WellKnownIdentifier int

const (
	IdentifierUnknown WellKnownIdentifier = iota

	IdentifierTestRootCA

	IdentifierIntCAService
	IdentifierIntCAIntranet
	IdentifierTestIntCA

	IdentifierDirectory

	DefaultPolicyIdCertRequest
	DefaultPolicyIdCertEnroll
	DefaultPolicyIdCertAadAppCredential
)

var (
	// nil is reserved

	// root CA --1 ~ --f
	WellKnownID_RootCA = uuid.MustParse("00000000-0000-0000-0000-000000000001")

	idTestRootCA = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-00000000000f"))

	// intermediate CAs --10 ~ --ff
	idIntCAService  = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-000000000011"))
	idIntCAIntranet = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-000000000012"))
	idTestIntCA     = WellKnownID(uuid.MustParse("00000000-0000-0000-0000-0000000000ff"))

	idDirectory = WellKnownID(uuid.MustParse(MustGetenv("AZURE_TENANT_ID")))

	// default policy ids --1-1 ~ --1-f
	defaultPolicyIdCertRequest          = WellKnownID(uuid.MustParse("00000000-0000-0000-0001-000000000001"))
	defaultPolicyIdCertEnroll           = WellKnownID(uuid.MustParse("00000000-0000-0000-0001-000000000002"))
	defaultPolicyIdCertAadAppCredential = WellKnownID(uuid.MustParse("00000000-0000-0000-0001-000000000003"))
)

var idMap = map[WellKnownIdentifier]WellKnownID{
	IdentifierTestRootCA: idTestRootCA,

	IdentifierIntCAService:  idIntCAService,
	IdentifierIntCAIntranet: idIntCAIntranet,
	IdentifierTestIntCA:     idTestIntCA,

	IdentifierDirectory: idDirectory,

	DefaultPolicyIdCertRequest:          defaultPolicyIdCertRequest,
	DefaultPolicyIdCertEnroll:           defaultPolicyIdCertEnroll,
	DefaultPolicyIdCertAadAppCredential: defaultPolicyIdCertAadAppCredential,
}

func GetID(identifier WellKnownIdentifier) uuid.UUID {
	return uuid.UUID(idMap[identifier])
}
