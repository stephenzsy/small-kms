package key

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

type KeyDoc struct {
	resdoc.ResourceDoc
	cloudkey.JsonWebKey
	Status        keymodels.KeyStatus  `json:"status"`
	Created       models.NumericDate   `json:"iat"`
	NotBefore     *models.NumericDate  `json:"nbf,omitempty"`
	NotAfter      *models.NumericDate  `json:"exp,omitempty"`
	Exportable    bool                 `json:"exportable"`
	Policy        resdoc.DocIdentifier `json:"policy"`
	PolicyVersion []byte               `json:"policyVersion"`
	Checksum      []byte               `json:"checksum"`
}

type keyGenerateDoc struct {
	KeyDoc

	rsaKeySize        int
	keyVaultStoreName string
}

func (d *keyGenerateDoc) init(nsProvider models.NamespaceProvider, nsID string, policy *KeyPolicyDoc) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.ID = id.String()
	d.KeyType = policy.KeySpec.Kty
	if policy.KeySpec.KeySize != nil {
		d.rsaKeySize = *policy.KeySpec.KeySize
	}
	d.Curve = policy.KeySpec.Crv
	d.KeyOperations = policy.KeySpec.KeyOperations
	d.Exportable = policy.Exportable
	if policy.ExpiryTime != nil {
		now := time.Now()
		d.NotBefore = jwt.NewNumericDate(now)
		d.NotAfter = jwt.NewNumericDate(caldur.Shift(now, *policy.ExpiryTime))
	}
	d.keyVaultStoreName = kv.GetMaterialName(kv.MaterialNameKindKey, nsProvider, nsID, policy.ID)
	d.Policy = policy.Identifier()
	d.PolicyVersion = policy.Version
	return nil
}

func (d *keyGenerateDoc) getAzCreateKeyParams() (params azkeys.CreateKeyParameters, err error) {
	switch d.KeyType {
	case cloudkey.KeyTypeEC:
		params.Kty = to.Ptr(azkeys.KeyTypeEC)
		switch d.Curve {
		case cloudkey.CurveNameP256:
			params.Curve = to.Ptr(azkeys.CurveNameP256)
		case cloudkey.CurveNameP384:
			params.Curve = to.Ptr(azkeys.CurveNameP384)
		case cloudkey.CurveNameP521:
			params.Curve = to.Ptr(azkeys.CurveNameP521)
		default:
			return params, cloudkey.ErrInvalidCurve
		}
	case cloudkey.KeyTypeRSA:
		params.Kty = to.Ptr(azkeys.KeyTypeRSA)
		switch d.rsaKeySize {
		case 2048, 3072, 4096:
			params.KeySize = to.Ptr(int32(d.rsaKeySize))
		}
	default:
		return params, cloudkey.ErrInvalidKeyType
	}
	// keyops
	params.KeyOps = make([]*azkeys.KeyOperation, len(d.KeyOperations))
	for i, keyOp := range d.KeyOperations {
		params.KeyOps[i] = to.Ptr(azkeys.KeyOperation(keyOp))
	}
	// exportable
	params.KeyAttributes = &azkeys.KeyAttributes{
		Exportable: &d.Exportable,
		Enabled:    to.Ptr(true),
	}
	if d.NotBefore != nil {
		params.KeyAttributes.NotBefore = &d.NotBefore.Time
	}
	if d.NotAfter != nil {
		params.KeyAttributes.Expires = &d.NotAfter.Time
	}
	return params, nil
}

func (d *KeyDoc) ToKeyRef() (m keymodels.KeyRef) {
	m.Ref = d.ToRef()
	m.Iat = d.Created
	m.Exp = d.NotAfter
	m.Status = d.Status
	m.PolicyIdentifier = d.Policy.String()
	return m
}

func (d *KeyDoc) ToModel() (m keymodels.Key) {
	m.KeyRef = d.ToKeyRef()
	m.Jwk = d.JsonWebKey
	m.Exportable = d.Exportable
	m.Identififier = d.Identifier().String()
	m.Nbf = d.NotBefore
	return m
}
