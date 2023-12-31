package base

import (
	"bytes"
	"encoding"
	"errors"
	"regexp"

	"github.com/google/uuid"
)

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

func (i identifier) String() string {
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

func StringIdentifier(id string) identifier {
	return identifier{
		isUUID: false,
		strVal: id,
	}
}

func ParseIdentifier(s string) (i identifier) {
	(&i).UnmarshalText([]byte(s))
	return
}

type NamespaceIdentifier struct {
	kind NamespaceKind
	id   ID
}

func (i NamespaceIdentifier) Kind() NamespaceKind {
	return i.kind
}

func (i NamespaceIdentifier) ID() ID {
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

func NewNamespaceIdentifier(nsKind NamespaceKind, nsID ID) NamespaceIdentifier {
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

func NewDocNamespacePartitionKey(nsKind NamespaceKind, nsID ID, rKind ResourceKind) DocNamespacePartitionKey {
	return DocNamespacePartitionKey{
		nsIdentifier: NewNamespaceIdentifier(nsKind, nsID),
		rKind:        rKind,
	}
}
