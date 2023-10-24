package base

import (
	"encoding"

	"github.com/google/uuid"
)

type identifierImpl struct {
	isUUID  bool
	uuidVal uuid.UUID
	strVal  string
}

// MarshalText implements encoding.TextMarshaler.
func (i identifierImpl) MarshalText() (text []byte, err error) {
	if i.isUUID {
		return i.uuidVal.MarshalText()
	}
	return []byte(i.strVal), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *identifierImpl) UnmarshalText(text []byte) error {
	*i = identifierImpl{}
	switch len(text) {
	case 32, 36, 36 + 2:
		var err error
		i.uuidVal, err = uuid.ParseBytes(text)
		if err == nil {
			i.isUUID = true
			return nil
		} else {
			i.uuidVal = uuid.UUID{}
		}
	}
	i.strVal = string(text)
	return nil
}

func (i *identifierImpl) String() string {
	if i == nil {
		return ""
	}
	if i.isUUID {
		return i.uuidVal.String()
	}
	return i.strVal
}

func (i *identifierImpl) IsUUID() bool {
	return i.isUUID
}

func (i *identifierImpl) UUID() uuid.UUID {
	return i.uuidVal
}

func (i *identifierImpl) Bytes() []byte {
	if i.isUUID {
		return i.uuidVal[:]
	}
	return []byte(i.strVal)
}

func (i *identifierImpl) AsUUID() (uuid.UUID, bool) {
	if !i.isUUID {
		return uuid.UUID{}, false
	}
	return i.uuidVal, i.isUUID
}

func (i *identifierImpl) IsNilOrEmpty() bool {
	if i == nil {
		return true
	}
	if i.isUUID {
		return i.uuidVal == uuid.Nil
	}
	return i.strVal == ""
}

func StringIdentifier[S ~string](s S) identifierImpl {
	return identifierImpl{
		isUUID: false,
		strVal: string(s),
	}
}

func UUIDIdentifier(uuid uuid.UUID) identifierImpl {
	return identifierImpl{
		isUUID:  true,
		uuidVal: uuid,
	}
}

func IdentifierFromString(s string) identifierImpl {
	return identifierImpl{
		isUUID: false,
		strVal: s,
	}
}

var _ encoding.TextMarshaler = identifierImpl{}
var _ encoding.TextUnmarshaler = (*identifierImpl)(nil)
