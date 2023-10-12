package common

import (
	"bytes"
	"encoding"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/shared"
)

var (
	ErrInvalidIdentifier = fmt.Errorf("invalid identifier")
)

// Deprecated: use shared.Identifier instead
type Identifier = shared.Identifier

// Deprecated: use shared.StringIdentifier instead
func StringIdentifier[T ~string](text T) Identifier {
	return shared.StringIdentifier[T](text)
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
