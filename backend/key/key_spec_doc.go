package key

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type KeySpecDoc struct {
	base.BaseDoc

	Name       string               `json:"name"`
	Exportable bool                 `json:"exportable"`
	KeyType    JsonWebKeyType       `json:"kty"`
	KeySize    *int32               `json:"key_size,omitempty"`
	CurveName  *JsonWebKeyCurveName `json:"crv,omitempty"`
	ExpiryTime *base.Period         `json:"expiryTime,omitempty"`
}

func NewKeySpecDoc(
	nsKind base.NamespaceKind,
	nsIdentifier base.Identifier,
	name string,
	params *KeySpecParameters,
) (*KeySpecDoc, error) {
	doc := &KeySpecDoc{
		BaseDoc: base.BaseDoc{
			NamespaceKind:       nsKind,
			NamespaceIdentifier: nsIdentifier,
			ResourceKind:        base.ResourceKindKeySpec,
			ResourceIdentifier:  base.StringIdentifier(name),
		},
		Name: name,
	}
	switch nsKind {
	case base.NamespaceKindRootCA:
		if params.Exportable != nil && *params.Exportable {
			return nil, fmt.Errorf("%w: exportable is not supported for root CA", base.ErrResponseStatusBadRequest)
		}
		doc.Exportable = false
	default:
		if params.Exportable != nil && *params.Exportable {
			doc.Exportable = true
		}
	}
	doc.KeyType = params.Kty
	switch doc.KeyType {
	case JsonWebKeyTypeEC:
		doc.CurveName = params.Crv
		if doc.CurveName == nil {
			doc.CurveName = utils.ToPtr(JsonWebKeyCurveNameP384)
		}
		switch *doc.CurveName {
		case JsonWebKeyCurveNameP256, JsonWebKeyCurveNameP384:
			// ok
		default:
			return nil, fmt.Errorf("%w: unsupported curve name: %s", base.ErrResponseStatusBadRequest, *doc.CurveName)
		}
	case JsonWebKeyTypeRSA:
		doc.KeySize = params.KeySize
		if doc.KeySize == nil {
			switch nsKind {
			case base.NamespaceKindRootCA:
				doc.KeySize = utils.ToPtr(int32(4096))
			default:
				doc.KeySize = utils.ToPtr(int32(2048))
			}
		}
		switch *doc.KeySize {
		case 2048, 3072, 4096:
			// ok
		default:
			return nil, fmt.Errorf("%w: unsupported key size: %d", base.ErrResponseStatusBadRequest, doc.KeySize)
		}
	default:
		return nil, fmt.Errorf("%w: unsupported key type: %s", base.ErrResponseStatusBadRequest, doc.KeyType)
	}
	doc.ExpiryTime = params.ExpiryTime
	return doc, nil
}
