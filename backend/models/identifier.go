package models

import "github.com/google/uuid"

type NameOrUUIDIdentifier struct {
	isUuid  bool
	uuidVal uuid.UUID
	strVal  string
}

func (identifier NameOrUUIDIdentifier) String() string {
	if identifier.isUuid {
		return identifier.uuidVal.String()
	}
	return identifier.strVal
}

func (identifier NameOrUUIDIdentifier) IsUUID() bool {
	return identifier.isUuid
}

func (identifier NameOrUUIDIdentifier) GetUUID() uuid.UUID {
	return identifier.uuidVal
}

func (identifier NameOrUUIDIdentifier) TryGetUUID() (uuid.UUID, bool) {
	return identifier.uuidVal, identifier.isUuid
}

func (identifier NameOrUUIDIdentifier) MarshalText() ([]byte, error) {
	if identifier.isUuid {
		return identifier.uuidVal.MarshalText()
	}
	return []byte(identifier.strVal), nil
}

func (identifier *NameOrUUIDIdentifier) IsNilOrEmpty() bool {
	if identifier.isUuid {
		return identifier.uuidVal == uuid.Nil
	}
	return identifier.strVal == ""
}

func (identifier *NameOrUUIDIdentifier) UnmarshalText(text []byte) (_ error) {
	// only accept length 36 or 38 for uuid
	*identifier = NameOrUUIDIdentifier{}
	l := len(text)
	if l == 36 || l == 38 {
		var err error
		if identifier.uuidVal, err = uuid.ParseBytes(text); err == nil {
			identifier.isUuid = true
			identifier.strVal = identifier.uuidVal.String()
			return
		} else {
			// clear
			identifier.uuidVal = uuid.UUID{}
		}
	}
	identifier.strVal = string(text)
	return
}

func IdentifierFromString(val string) NameOrUUIDIdentifier {
	var identifier NameOrUUIDIdentifier
	identifier.UnmarshalText([]byte(val))
	return identifier
}
