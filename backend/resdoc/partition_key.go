package resdoc

import (
	"encoding"
	"errors"
	"strings"

	"github.com/stephenzsy/small-kms/backend/models"
)

var (
	ErrInvalidPartitionKey = errors.New("invalid partition key")
)

type PartitionKey struct {
	NamespaceProvider models.NamespaceProvider
	NamespaceID       string
	ResourceProvider  models.ResourceProvider
}

func (p PartitionKey) String() string {
	return string(p.NamespaceProvider) + ":" + p.NamespaceID + ":" + string(p.ResourceProvider)
}

func (p PartitionKey) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *PartitionKey) UnmarshalText(text []byte) (err error) {
	*p, err = PartitionKeyFromString(string(text))
	return err
}

func PartitionKeyFromString(text string) (PartitionKey, error) {
	var p PartitionKey
	if text == "" {
		return p, nil
	}
	parts := strings.Split(text, ":")
	if len(parts) == 3 {
		p.NamespaceProvider = models.NamespaceProvider(parts[0])
		p.NamespaceID = parts[1]
		p.ResourceProvider = models.ResourceProvider(parts[2])
		return p, nil
	}
	return p, ErrInvalidPartitionKey
}

var _ encoding.TextMarshaler = PartitionKey{}
var _ encoding.TextUnmarshaler = (*PartitionKey)(nil)

type DocIdentifier struct {
	PartitionKey
	ID string
}
