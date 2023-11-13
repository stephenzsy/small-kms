package key

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type KeyDoc struct {
	base.BaseDoc
	cloudkey.JsonWebKeyBase
	KeySize    *int32            `json:"keySize"`
	NotAfter   *base.NumericDate `json:"exp"`
	Exportable bool              `json:"exportable"`
}

func (d *KeyDoc) init(c context.Context, policy *KeyPolicyDoc) error {
	nsCtx := ns.GetNSContext(c)
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.BaseDoc.Init(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindKey, base.IDFromUUID(id))
	return nil
}

func (d *KeyDoc) toAzCreateKeyParameters() azkeys.CreateKeyParameters {

	p := azkeys.CreateKeyParameters{
		KeyAttributes: &azkeys.KeyAttributes{
			Enabled:    to.Ptr(true),
			Exportable: to.Ptr(d.Exportable),
		},
	}
	switch d.KeyType {
	case cloudkey.KeyTypeRSA:
		p.Kty = to.Ptr(azkeys.KeyTypeRSA)
		p.KeySize = d.KeySize
	case cloudkey.KeyTypeEC:
		p.Kty = to.Ptr(azkeys.KeyTypeEC)
		switch d.Curve {
		case cloudkey.CurveNameP256:
			p.Curve = to.Ptr(azkeys.CurveNameP256)
		case cloudkey.CurveNameP384:
			p.Curve = to.Ptr(azkeys.CurveNameP384)
		case cloudkey.CurveNameP521:
			p.Curve = to.Ptr(azkeys.CurveNameP521)
		}
	case cloudkey.KeyTypeOct:
		p.Kty = to.Ptr(azkeys.KeyTypeOct)
	}

	p.KeyOps = utils.MapSlice(d.KeyOperations, func(keyOp JsonWebKeyOperation) *azkeys.KeyOperation {
		return to.Ptr(cloudkeyaz.ToAzKeyOperation(keyOp))
	})

	if d.NotAfter != nil {
		p.KeyAttributes.Expires = &d.NotAfter.Time

	}

	return p
}
