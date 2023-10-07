package common

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
)

var (
	ErrInvalidLocator = errors.New("invalid locator")
)

type Locator[NK ~string, K ~string] struct {
	nsID IdentifierWithKind[NK]
	id   IdentifierWithKind[K]
}

func (l Locator[NK, K]) GetNamespaceID() IdentifierWithKind[NK] {
	return l.nsID
}

func (l Locator[NK, K]) GetID() IdentifierWithKind[K] {
	return l.id
}

func (l Locator[NK, K]) String() string {
	return fmt.Sprintf("%s/%s", l.nsID, l.id)
}

// MarshalText implements encoding.TextMarshaler.
func (l Locator[NK, K]) MarshalText() ([]byte, error) {
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
func (l *Locator[NK, K]) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte("/"), 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w:%s", ErrInvalidLocator, text)
	}
	l = &Locator[NK, K]{}
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

func NewLocator[NK ~string, K ~string](namespaceID IdentifierWithKind[NK], resourceID IdentifierWithKind[K]) Locator[NK, K] {
	return Locator[NK, K]{
		nsID: namespaceID,
		id:   resourceID,
	}
}

var _ encoding.TextMarshaler = Locator[string, string]{}
var _ encoding.TextUnmarshaler = (*Locator[string, string])(nil)
