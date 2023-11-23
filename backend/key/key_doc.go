package key

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type KeyDoc struct {
	base.BaseDoc
	cloudkey.JsonWebKey
	KeySize       *int32            `json:"keySize,omitempty"`
	Created       base.NumericDate  `json:"iat"`
	NotBefore     *base.NumericDate `json:"nbf,omitempty"`
	NotAfter      *base.NumericDate `json:"exp,omitempty"`
	Exportable    bool              `json:"exportable"`
	Policy        base.DocLocator   `json:"policy"`
	PolicyVersion base.HexDigest    `json:"policyVersion"`
}

func (d *KeyDoc) init(c context.Context, policy *KeyPolicyDoc) error {
	nsCtx := ns.GetNSContext(c)
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.BaseDoc.Init(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindKey, base.IDFromUUID(id))
	d.KeyType = policy.KeyProperties.Kty
	d.KeySize = policy.KeyProperties.KeySize
	d.Curve = policy.KeyProperties.Crv
	d.KeyOperations = policy.KeyProperties.KeyOperations
	d.Exportable = policy.Exportable
	if policy.ExpiryTime != nil {
		d.NotAfter = jwt.NewNumericDate(base.AddPeriod(time.Now(), *policy.ExpiryTime))
	}
	d.Policy = policy.GetStorageFullIdentifier()
	d.PolicyVersion = policy.Version
	return nil
}

func (d *KeyDoc) populateModelRef(r *KeyRef) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&r.ResourceReference)
	r.Iat = d.Created
	r.Exp = d.NotAfter
}

func (d *KeyDoc) populateModel(r *Key) {
	if d == nil || r == nil {
		return
	}
	d.populateModelRef(&r.KeyRef)
	r.KeyType = d.KeyType
	r.KeySize = d.KeySize
	r.Curve = d.Curve
	r.N = d.N
	r.E = d.E
	r.X = d.X
	r.Y = d.Y
	r.Nbf = d.NotBefore
	r.Exp = d.NotAfter
	r.Iat = d.Created
	r.KeyID = d.KeyID
	r.Policy = d.Policy
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
