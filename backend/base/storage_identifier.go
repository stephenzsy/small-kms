package base

import (
	"bytes"
	"encoding"
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

type identifier struct {
	isUUID  bool
	uuidVal UUID
	strVal  string
}

var identifierRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *identifier) UnmarshalText(text []byte) error {
	if len(text) == 36 || len(text) == 36+2 {
		if parsedUUID, err := uuid.ParseBytes(text); err == nil {
			*p = identifier{
				isUUID:  true,
				uuidVal: parsedUUID,
			}
			return nil
		}
	}
	*p = identifier{
		isUUID: false,
		strVal: string(text),
	}
	if p.strVal == "" {
		p.strVal = "default"
	} else if len(p.strVal) > 48 {
		return errors.New("identifier too long, max 48 characters")
	} else if len(p.strVal) < 2 {
		return errors.New("identifier too short min 2 characters")
	} else if !identifierRegex.MatchString(p.strVal) {
		return errors.New("invalid identifier format: " + p.strVal)
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (k identifier) MarshalText() (text []byte, err error) {
	if k.isUUID {
		return k.uuidVal.MarshalText()
	}
	return []byte(k.strVal), nil
}

func (i Identifier) String() string {
	if i.isUUID {
		return i.uuidVal.String()
	}
	return i.strVal
}

func (k identifier) Length() int {
	if k.isUUID {
		return 36
	}
	return len(k.strVal)
}

func (k identifier) IsUUID() bool {
	return k.isUUID
}

func (k identifier) UUID() UUID {
	return k.uuidVal
}

func (k *identifier) IsNilOrEmpty() bool {
	if k == nil {
		return true
	}
	if k.isUUID {
		return k.uuidVal == uuid.Nil
	}
	return k.strVal == ""
}

var _ encoding.TextMarshaler = identifier{}
var _ encoding.TextUnmarshaler = (*identifier)(nil)

func UUIDIdentifier(id UUID) Identifier {
	return Identifier{
		isUUID:  true,
		uuidVal: id,
	}
}

func StringIdentifier(id string) Identifier {
	return Identifier{
		isUUID: false,
		strVal: id,
	}
}

func ParseIdentifier(s string) (i Identifier) {
	(&i).UnmarshalText([]byte(s))
	return
}

type NamespaceIdentifier struct {
	kind NamespaceKind
	id   Identifier
}

func (i NamespaceIdentifier) Kind() NamespaceKind {
	return i.kind
}

func (i NamespaceIdentifier) Identifier() Identifier {
	return i.id
}

func (i NamespaceIdentifier) String() string {
	if b, err := i.MarshalText(); err != nil {
		return ""
	} else {
		return string(b)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *NamespaceIdentifier) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte(":"), 2)
	if len(parts) != 2 {
		return errors.New("invalid storage namespace ID format")
	}
	*p = NamespaceIdentifier{
		kind: NamespaceKind(parts[0]),
	}
	return p.id.UnmarshalText(parts[1])
}

// MarshalText implements encoding.TextMarshaler.
func (i NamespaceIdentifier) MarshalText() (text []byte, err error) {
	if keyBytes, err := i.id.MarshalText(); err != nil {
		return nil, err
	} else {
		text = make([]byte, len(i.kind)+len(keyBytes)+1)
		copy(text, i.kind)
		text[len(i.kind)] = ':'
		copy(text[len(i.kind)+1:], keyBytes)
	}
	return
}

var _ encoding.TextMarshaler = NamespaceIdentifier{}
var _ encoding.TextUnmarshaler = (*NamespaceIdentifier)(nil)

func NewNamespaceIdentifier(nsKind NamespaceKind, nsID identifier) NamespaceIdentifier {
	return NamespaceIdentifier{
		kind: nsKind,
		id:   nsID,
	}
}

type DocNamespacePartitionKey struct {
	nsIdentifier NamespaceIdentifier
	rKind        ResourceKind
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *DocNamespacePartitionKey) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte(":"), 3)
	if len(parts) != 3 {
		return errors.New("invalid storage namespace ID format")
	}
	*p = DocNamespacePartitionKey{
		nsIdentifier: NamespaceIdentifier{
			kind: NamespaceKind(parts[0]),
		},
		rKind: ResourceKind(parts[2]),
	}
	return p.nsIdentifier.id.UnmarshalText(parts[1])
}

// MarshalText implements encoding.TextMarshaler.
func (i DocNamespacePartitionKey) MarshalText() (text []byte, err error) {
	if nsIDBytes, err := i.nsIdentifier.MarshalText(); err != nil {
		return nil, err
	} else {
		text = make([]byte, len(nsIDBytes)+len(i.rKind)+1)
		copy(text, nsIDBytes)
		text[len(nsIDBytes)] = ':'
		copy(text[len(nsIDBytes)+1:], i.rKind)
	}
	return
}

func (i DocNamespacePartitionKey) String() string {
	if b, err := i.MarshalText(); err != nil {
		return ""
	} else {
		return string(b)
	}
}

var _ encoding.TextMarshaler = DocNamespacePartitionKey{}
var _ encoding.TextUnmarshaler = (*DocNamespacePartitionKey)(nil)

func NewDocNamespacePartitionKey(nsKind NamespaceKind, nsID Identifier, rKind ResourceKind) DocNamespacePartitionKey {
	return DocNamespacePartitionKey{
		nsIdentifier: NewNamespaceIdentifier(nsKind, nsID),
		rKind:        rKind,
	}
}

type DocFullIdentifier struct {
	pKey  DocNamespacePartitionKey
	docID identifier
}

func (i DocFullIdentifier) PartitionKey() DocNamespacePartitionKey {
	return i.pKey
}

func (i DocFullIdentifier) NamespaceKind() NamespaceKind {
	return i.pKey.nsIdentifier.kind
}

func (i DocFullIdentifier) NamespaceIdentifier() Identifier {
	return i.pKey.nsIdentifier.id
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *DocFullIdentifier) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte("/"), 2)
	if len(parts) != 2 {
		return errors.New("invalid identifier format")
	}
	*p = DocFullIdentifier{}
	if err := p.pKey.UnmarshalText(parts[0]); err != nil {
		return err
	}
	return p.docID.UnmarshalText(parts[1])
}

// MarshalText implements encoding.TextMarshaler.
func (i DocFullIdentifier) MarshalText() (text []byte, err error) {
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

var _ encoding.TextMarshaler = DocFullIdentifier{}
var _ encoding.TextUnmarshaler = (*DocFullIdentifier)(nil)

func (i DocFullIdentifier) String() string {
	if b, err := i.MarshalText(); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func NewDocFullIdentifier(nsKind NamespaceKind, nsID Identifier, rKind ResourceKind, rID Identifier) DocFullIdentifier {
	return DocFullIdentifier{
		pKey:  NewDocNamespacePartitionKey(nsKind, nsID, rKind),
		docID: rID,
	}
}
