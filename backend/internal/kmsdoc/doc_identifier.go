package kmsdoc

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/stephenzsy/small-kms/backend/common"
)

var ErrDocIdentifierInvalid = fmt.Errorf("invalid doc identifier")

// document name space type
type DocNsType string

const (
	DocNsTypeProfile   DocNsType = "profile"
	DocNsTypeCaRoot    DocNsType = "ca-root"
	DocNsTypeCaInt     DocNsType = "ca-int"
	DocNSTypeDirectory DocNsType = "directory"
)

type DocKind string

const (
	DocKindCaRoot              DocKind = "ca-root"       // only for profile/builtin
	DocKindCaInt               DocKind = "ca-int"        // only for profile/builtin
	DocKindDirectoryObject     DocKind = "object"        // only for profile/tenant
	DocKindCertificateTemplate DocKind = "cert-template" //
)

type DocIdentifier[T DocNsType | DocKind] struct {
	kind       T
	identifier common.Identifier
}

func (d DocIdentifier[T]) String() string {
	return string(d.kind) + "/" + d.identifier.String()
}

func (d DocIdentifier[T]) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d DocIdentifier[T]) Kind() T {
	return d.kind
}

func (d DocIdentifier[T]) Identifier() common.Identifier {
	return d.identifier
}

func (d *DocIdentifier[T]) UnmarshalText(text []byte) (err error) {
	id := string(text)
	l := strings.SplitN(id, "/", 2)
	if len(l) != 2 {
		return fmt.Errorf("%w:%s", ErrDocIdentifierInvalid, string(text))
	}
	d.kind = T(l[0])
	d.identifier = common.StringIdentifier(l[1])
	return
}

func NewDocIdentifier[T DocNsType | DocKind](t T, id common.Identifier) DocIdentifier[T] {
	return DocIdentifier[T]{kind: t, identifier: id}
}

var _ encoding.TextMarshaler = &DocIdentifier[DocNsType]{}
var _ encoding.TextUnmarshaler = &DocIdentifier[DocNsType]{}

type DocNsID = DocIdentifier[DocNsType]
type DocID = DocIdentifier[DocKind]

func StringDocIdentifier[T DocNsType | DocKind](nsType T, id string) DocIdentifier[T] {
	return DocIdentifier[T]{kind: nsType, identifier: common.StringIdentifier(id)}
}
