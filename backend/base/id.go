package base

import (
	"encoding"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

type ID string

// MarshalText implements encoding.TextMarshaler.
func (id ID) MarshalText() (text []byte, err error) {
	return []byte(id), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (idPtr *ID) UnmarshalText(text []byte) error {
	if parsed, err := uuid.ParseBytes(text); err != nil {
		*idPtr = ID(string(text))
	} else {
		*idPtr = ID(parsed.String())
	}
	return nil
}

var _ encoding.TextMarshaler = ID("")
var _ encoding.TextUnmarshaler = (*ID)(nil)

func (id ID) UUID() UUID {
	uuid, _ := uuid.Parse(string(id))
	return uuid
}

func (id ID) AsUUID() (UUID, bool) {
	uuid, err := uuid.Parse(string(id))
	return uuid, err == nil
}

func IDFromString(s string) ID {
	return ID(s)
}

func IDFromUUID(u UUID) ID {
	return ID(u.String())
}

func ParseID(s string) (id ID) {
	id.UnmarshalText([]byte(s))
	return id
}
