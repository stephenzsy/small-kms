package shared

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

type identifierImpl struct {
	isUuid  bool
	uuidVal uuid.UUID
	strVal  string
}

func (identifier identifierImpl) String() string {
	if identifier.isUuid {
		return identifier.uuidVal.String()
	}
	return identifier.strVal
}

func (identifier identifierImpl) IsUUID() bool {
	return identifier.isUuid
}

func (identifier identifierImpl) UUID() uuid.UUID {
	return identifier.uuidVal
}

func (identifier identifierImpl) UUIDPtr() *uuid.UUID {
	return &identifier.uuidVal
}

func (identifier identifierImpl) TryGetUUID() (uuid.UUID, bool) {
	return identifier.uuidVal, identifier.isUuid
}

func (identifier identifierImpl) MarshalText() ([]byte, error) {
	if identifier.isUuid {
		return identifier.uuidVal.MarshalText()
	}
	return []byte(identifier.strVal), nil
}

func (identifier *identifierImpl) IsNilOrEmpty() bool {
	if identifier == nil {
		return true
	}
	if identifier.isUuid {
		return identifier.uuidVal == uuid.Nil
	}
	return identifier.strVal == ""
}

func (identifier *identifierImpl) UnmarshalText(text []byte) (_ error) {
	*identifier = identifierFromTextBytes(text)
	return
}

func identifierFromTextBytes(text []byte) (identifier identifierImpl) {
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
func StringIdentifier[T ~string](text T) identifierImpl {
	return identifierImpl{
		isUuid: false,
		strVal: string(text),
	}
}

func UUIDIdentifier(uuid uuid.UUID) identifierImpl {
	return identifierImpl{
		isUuid:  true,
		uuidVal: uuid,
		strVal:  uuid.String(),
	}
}

func UUIDIdentifierFromString(text string) identifierImpl {
	return identifierFromTextBytes([]byte(text))
}

func UUIDIdentifierFromStringPtr(p *string) identifierImpl {
	if p == nil {
		return identifierImpl{}
	}
	return identifierFromTextBytes([]byte(*p))
}

var _ encoding.TextMarshaler = &identifierImpl{}
var _ encoding.TextUnmarshaler = &identifierImpl{}

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

func (identifier identifierImpl) HasReservedIDOrPrefix() bool {
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

func (identifier identifierImpl) IsValid() bool {
	return identifier.isUuid || (len(identifier.strVal) < 128 && identifierRegex.MatchString(identifier.strVal))
}

// immutable externally
type identifierWithKind[K ~string] struct {
	kind       K
	identifier identifierImpl
}

func (d identifierWithKind[K]) String() string {
	return string(d.kind) + ":" + d.identifier.String()
}

func (d identifierWithKind[K]) Kind() K {
	return d.kind
}

func (d identifierWithKind[K]) Identifier() identifierImpl {
	return d.identifier
}

func (d identifierWithKind[K]) MarshalText() (text []byte, err error) {
	kindBytes := []byte(d.kind)
	identifierBytes, err := d.identifier.MarshalText()
	return bytes.Join([][]byte{kindBytes, identifierBytes}, []byte(":")), err
}

func (d *identifierWithKind[K]) UnmarshalText(text []byte) (err error) {
	l := bytes.SplitN(text, []byte(":"), 2)
	if len(l) != 2 {
		return fmt.Errorf("%w:%s", ErrInvalidIdentifier, text)
	}
	*d = identifierWithKind[K]{
		kind: K(l[0]),
	}
	d.identifier.UnmarshalText(l[1])
	return
}

func (c identifierWithKind[K]) WithKind(kind K) identifierWithKind[K] {
	return identifierWithKind[K]{
		kind:       kind,
		identifier: c.identifier,
	}
}

var _ encoding.TextMarshaler = identifierWithKind[string]{}
var _ encoding.TextUnmarshaler = (*identifierWithKind[string])(nil)

func NewNamespaceIdentifier(kind NamespaceKind, identifier Identifier) identifierWithKind[NamespaceKind] {
	return identifierWithKind[NamespaceKind]{
		kind:       kind,
		identifier: identifier,
	}
}

func NewResourceIdentifier(kind ResourceKind, identifier Identifier) identifierWithKind[ResourceKind] {
	return identifierWithKind[ResourceKind]{
		kind:       kind,
		identifier: identifier,
	}
}
