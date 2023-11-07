package base

import (
	"bytes"
	"encoding"
	"errors"
)

type DocLocator struct {
	pKey  DocNamespacePartitionKey
	docID ID
}

func (i DocLocator) PartitionKey() DocNamespacePartitionKey {
	return i.pKey
}

func (i DocLocator) NamespaceKind() NamespaceKind {
	return i.pKey.nsIdentifier.kind
}

func (i DocLocator) NamespaceID() ID {
	return i.pKey.nsIdentifier.id
}

func (i DocLocator) ID() ID {
	return i.docID
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *DocLocator) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte("/"), 2)
	if len(parts) != 2 {
		return errors.New("invalid identifier format")
	}
	*p = DocLocator{}
	if err := p.pKey.UnmarshalText(parts[0]); err != nil {
		return err
	}
	return p.docID.UnmarshalText(parts[1])
}

// MarshalText implements encoding.TextMarshaler.
func (i DocLocator) MarshalText() (text []byte, err error) {
	if nsIDBytes, err := i.pKey.MarshalText(); err != nil {
		return nil, err
	} else if docIDBytes, err := i.docID.MarshalText(); err != nil {
		return nil, err
	} else {
		text = make([]byte, len(nsIDBytes)+len(docIDBytes)+1)
		copy(text, nsIDBytes)
		text[len(nsIDBytes)] = '/'
		copy(text[len(nsIDBytes)+1:], docIDBytes)
	}
	return
}

var _ encoding.TextMarshaler = DocLocator{}
var _ encoding.TextUnmarshaler = (*DocLocator)(nil)

func (i DocLocator) String() string {
	if b, err := i.MarshalText(); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func NewDocLocator(nsKind NamespaceKind, nsID ID, rKind ResourceKind, rID ID) DocLocator {
	return DocLocator{
		pKey:  NewDocNamespacePartitionKey(nsKind, nsID, rKind),
		docID: rID,
	}
}
