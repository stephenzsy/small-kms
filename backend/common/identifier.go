package common

import (
	"encoding"

	"github.com/google/uuid"
)

type Identifier struct {
	isUuid  bool
	uuidVal uuid.UUID
	strVal  string
}

func (identifier Identifier) String() string {
	if identifier.isUuid {
		return identifier.uuidVal.String()
	}
	return identifier.strVal
}

func (identifier Identifier) IsUUID() bool {
	return identifier.isUuid
}

func (identifier Identifier) GetUUID() uuid.UUID {
	return identifier.uuidVal
}

func (identifier Identifier) TryGetUUID() (uuid.UUID, bool) {
	return identifier.uuidVal, identifier.isUuid
}

func (identifier Identifier) MarshalText() ([]byte, error) {
	if identifier.isUuid {
		return identifier.uuidVal.MarshalText()
	}
	return []byte(identifier.strVal), nil
}

func (identifier *Identifier) IsEmpty() bool {
	if identifier.isUuid {
		return identifier.uuidVal == uuid.Nil
	}
	return identifier.strVal == ""
}

func (identifier *Identifier) UnmarshalText(text []byte) (_ error) {
	*identifier = identifierFromTextBytes(text)
	return
}

func identifierFromTextBytes(text []byte) (identifier Identifier) {
	l := len(text)
	// only accept length 36 or 38 for uuid
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

// construct identifier from string, MUST NOT use possible uuid string or as will lead into type consistency
func StringIdentifier(text string) Identifier {
	return Identifier{
		isUuid: false,
		strVal: text,
	}
}

func UUIDIdentifier(uuid uuid.UUID) Identifier {
	return Identifier{
		isUuid:  true,
		uuidVal: uuid,
		strVal:  uuid.String(),
	}
}

func UUIDIdentifierFromString(text string) Identifier {
	return identifierFromTextBytes([]byte(text))
}

func UUIDIdentifierFromStringPtr(p *string) Identifier {
	if p == nil {
		return Identifier{}
	}
	return identifierFromTextBytes([]byte(*p))
}

var _ encoding.TextMarshaler = &Identifier{}
var _ encoding.TextUnmarshaler = &Identifier{}
