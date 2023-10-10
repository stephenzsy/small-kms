package common

import (
	"bytes"
	"encoding"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidIdentifier = fmt.Errorf("invalid identifier")
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

func (identifier Identifier) UUID() uuid.UUID {
	return identifier.uuidVal
}

func (identifier Identifier) UUIDPtr() *uuid.UUID {
	return &identifier.uuidVal
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

func (identifier *Identifier) IsNilOrEmpty() bool {
	if identifier == nil {
		return true
	}
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
func StringIdentifier[T ~string](text T) Identifier {
	return Identifier{
		isUuid: false,
		strVal: string(text),
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

const (
	reservedPrefixDefault  = "default"
	reservedPrefixSystem   = "system"
	reservedPrefixReserved = "reserved"
)

var reservedPrefixes map[string]bool = map[string]bool{
	reservedPrefixDefault:  true,
	reservedPrefixSystem:   true,
	reservedPrefixReserved: true,
	"latest":               true,
	"template":             true,
	"self":                 true,
	"global":               true,
}

func (identifier Identifier) HasReservedIDOrPrefix() bool {
	if identifier.isUuid {
		return identifier.uuidVal.Version() != 4
	}
	s := identifier.strVal
	lenStr := len(s)
	if lenStr <= 3 {
		return false
	}
	return reservedPrefixes[strings.ToLower(s)]
}

var identifierRegex = regexp.MustCompile("[A-Za-z0-9_-]+")

func (identifier Identifier) IsValid() bool {
	return identifier.isUuid || (len(identifier.strVal) < 128 && identifierRegex.MatchString(identifier.strVal))
}

// immutable externally
type IdentifierWithKind[K ~string] struct {
	kind       K
	identifier Identifier
}

func (d IdentifierWithKind[K]) String() string {
	return string(d.kind) + ":" + d.identifier.String()
}

func (d IdentifierWithKind[K]) Kind() K {
	return d.kind
}

func (d IdentifierWithKind[K]) Identifier() Identifier {
	return d.identifier
}

func (d IdentifierWithKind[K]) MarshalText() (text []byte, err error) {
	kindBytes := []byte(d.kind)
	identifierBytes, err := d.identifier.MarshalText()
	return bytes.Join([][]byte{kindBytes, identifierBytes}, []byte(":")), err
}

func (d *IdentifierWithKind[K]) UnmarshalText(text []byte) (err error) {
	l := bytes.SplitN(text, []byte(":"), 2)
	if len(l) != 2 {
		return fmt.Errorf("%w:%s", ErrInvalidIdentifier, text)
	}
	*d = IdentifierWithKind[K]{
		kind: K(l[0]),
	}
	d.identifier.UnmarshalText(l[1])
	return
}

func (c IdentifierWithKind[K]) WithKind(kind K) IdentifierWithKind[K] {
	return IdentifierWithKind[K]{
		kind:       kind,
		identifier: c.identifier,
	}
}

func NewIdentifierWithKind[K ~string](kind K, identifier Identifier) IdentifierWithKind[K] {
	return IdentifierWithKind[K]{kind: kind, identifier: identifier}
}

var _ encoding.TextMarshaler = IdentifierWithKind[string]{}
var _ encoding.TextUnmarshaler = (*IdentifierWithKind[string])(nil)
