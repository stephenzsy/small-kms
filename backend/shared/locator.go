package shared

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
)

var (
	ErrInvalidLocator = errors.New("invalid locator")
)

type locator[NK ~string, K ~string] struct {
	nsID identifierWithKind[NK]
	id   identifierWithKind[K]
}

func (l locator[NK, K]) GetNamespaceID() identifierWithKind[NK] {
	return l.nsID
}

func (l locator[NK, K]) GetID() identifierWithKind[K] {
	return l.id
}

func (l locator[NK, K]) String() string {
	return fmt.Sprintf("%s/%s", l.nsID, l.id)
}

func (l *locator[NK, K]) IsNilOrEmpty() bool {
	if l == nil {
		return true
	}
	return l.nsID.IsNilEmpty() && l.id.IsNilEmpty()
}

// MarshalText implements encoding.TextMarshaler.
func (l locator[NK, K]) MarshalText() ([]byte, error) {
	nsBytes, err := l.nsID.MarshalText()
	if err != nil {
		return nil, err
	}
	idBytes, err := l.id.MarshalText()
	if err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{nsBytes, idBytes}, []byte("/")), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (l *locator[NK, K]) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte("/"), 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w:%s", ErrInvalidLocator, text)
	}
	*l = locator[NK, K]{}
	err := l.nsID.UnmarshalText(parts[0])
	if err != nil {
		return fmt.Errorf("%w:%w", ErrInvalidLocator, err)
	}
	err = l.id.UnmarshalText(parts[1])
	if err != nil {
		return fmt.Errorf("%w:%w", ErrInvalidLocator, err)
	}
	return nil
}

func newLocator[NK ~string, K ~string](namespaceID identifierWithKind[NK], resourceID identifierWithKind[K]) locator[NK, K] {
	return locator[NK, K]{
		nsID: namespaceID,
		id:   resourceID,
	}
}

func (c locator[NK, K]) WithID(id identifierWithKind[K]) locator[NK, K] {
	return locator[NK, K]{
		nsID: c.nsID,
		id:   id,
	}
}

func (c locator[NK, K]) WithIDKind(kind K) locator[NK, K] {
	return locator[NK, K]{
		nsID: c.nsID,
		id:   c.id.WithKind(kind),
	}
}

var _ encoding.TextMarshaler = locator[string, string]{}
var _ encoding.TextUnmarshaler = (*locator[string, string])(nil)

func NewResourceLocator(namespaceID NamespaceIdentifier, resourceID ResourceIdentifier) ResourceLocator {
	return newLocator(namespaceID, resourceID)
}
