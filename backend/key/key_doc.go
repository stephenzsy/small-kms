package key

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
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
