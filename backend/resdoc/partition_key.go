package resdoc

import (
	"encoding"
	"errors"
	"strings"

	"github.com/stephenzsy/small-kms/backend/models"
)

var (
	ErrInvalidPartitionKey = errors.New("invalid partition key")
	ErrInvalidIdentifier   = errors.New("invalid identifier")
)

type PartitionKey struct {
	NamespaceProvider models.NamespaceProvider
	NamespaceID       string
	ResourceProvider  models.ResourceProvider
}

func (p PartitionKey) String() string {
	if p.NamespaceProvider == "" && p.NamespaceID == "" && p.ResourceProvider == "" {
		return ""
	}
	return string(p.NamespaceProvider) + ":" + p.NamespaceID + ":" + string(p.ResourceProvider)
}

func (p PartitionKey) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *PartitionKey) UnmarshalText(text []byte) (err error) {
	*p, err = ParsePartitionKey(string(text))
	return err
}

func (p *PartitionKey) IsEmpty() bool {
	return p.NamespaceProvider == "" && p.NamespaceID == "" && p.ResourceProvider == ""
}

func ParsePartitionKey(text string) (PartitionKey, error) {
	var p PartitionKey
	if text == "" {
		return p, nil
	}
	parts := strings.SplitN(text, ":", 3)
	if len(parts) != 3 {
		return p, ErrInvalidPartitionKey
	}
	p.NamespaceProvider = models.NamespaceProvider(parts[0])
	p.NamespaceID = parts[1]
	p.ResourceProvider = models.ResourceProvider(parts[2])
	return p, nil
}

var _ encoding.TextMarshaler = PartitionKey{}
var _ encoding.TextUnmarshaler = (*PartitionKey)(nil)

type DocIdentifier struct {
	PartitionKey
	ID string
}

func (p DocIdentifier) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *DocIdentifier) UnmarshalText(text []byte) (err error) {
	*p, err = ParseIdentifier(string(text))
	return err
}

func ParseIdentifier(text string) (identifier DocIdentifier, err error) {
	if text == "" {
		return identifier, nil
	}
	parts := strings.SplitN(text, "/", 2)
	if len(parts) != 2 {
		return identifier, ErrInvalidIdentifier
	}
	identifier.PartitionKey, err = ParsePartitionKey(parts[0])
	if err != nil {
		return identifier, err
	}
	identifier.ID = parts[1]
	return identifier, nil
}

func (p DocIdentifier) String() string {
	partitionKeyString := p.PartitionKey.String()
	if partitionKeyString == "" && p.ID == "" {
		return ""
	}
	return partitionKeyString + "/" + p.ID
}

func (p *DocIdentifier) IsEmpty() bool {
	return p.PartitionKey.IsEmpty() && p.ID == ""
}

var _ encoding.TextMarshaler = DocIdentifier{}
var _ encoding.TextUnmarshaler = (*DocIdentifier)(nil)
