package key

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type KeyPolicyDoc struct {
	base.BaseDoc
	DisplayName   string                       `json:"displayName"`
	KeyProperties GenerateJsonWebKeyProperties `json:"keyProperties"`
	Exportable    bool                         `json:"exportable"`
	Version       base.HexDigest               `json:"version"`
	ExpiryTime    *base.Period                 `json:"expiryTime"`
}

const (
	queryColumnDisplayName = "c.displayName"
)

func (doc *KeyPolicyDoc) init(c ctx.RequestContext, policyID base.ID, params *KeyPolicyParameters) error {
	logger := log.Ctx(c)
	nsCtx := ns.GetNSContext(c)
	doc.BaseDoc.Init(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindKeyPolicy, policyID)

	if params.DisplayName != "" {
		doc.DisplayName = params.DisplayName
	} else {
		doc.DisplayName = string(policyID)
	}

	digester := md5.New()
	// key properties default
	doc.KeyProperties = GenerateJsonWebKeyProperties{
		Kty:     cloudkey.KeyTypeRSA,
		KeySize: utils.ToPtr(int32(2048)),
		KeyOperations: []cloudkey.JsonWebKeyOperation{
			cloudkey.JsonWebKeyOperationSign,
			cloudkey.JsonWebKeyOperationVerify,
		},
	}
	if params.KeyProperties != nil {
		switch params.KeyProperties.Kty {
		case cloudkey.KeyTypeRSA:
			doc.KeyProperties.Kty = cloudkey.KeyTypeRSA
			doc.KeyProperties.KeySize = utils.ToPtr(int32(2048))
			if params.KeyProperties.KeySize != nil {
				switch *params.KeyProperties.KeySize {
				case 2048, 3072, 4096:
					doc.KeyProperties.KeySize = params.KeyProperties.KeySize
				case 0:
					// use default
				default:
					logger.Warn().Int32("keySize", *params.KeyProperties.KeySize).Msg("invalid key size, default to 2048")
				}
			}

		case cloudkey.KeyTypeEC:
			doc.KeyProperties.Kty = cloudkey.KeyTypeEC
			doc.KeyProperties.Crv = cloudkey.CurveNameP384 // default to P384 as cost in KeyVault is the same as 256

			switch params.KeyProperties.Crv {
			case cloudkey.CurveNameP256, cloudkey.CurveNameP384, cloudkey.CurveNameP521:
				doc.KeyProperties.Crv = params.KeyProperties.Crv
			case "":
				// default
			default:
				logger.Warn().Str("crv", string(params.KeyProperties.Crv)).Msg("invalid curve, default to P384")
			}
		}
	}

	doc.KeyProperties.writeToDigest(digester)

	if params.Exportable != nil && *params.Exportable {
		doc.Exportable = true
		digester.Write([]byte("exportable"))
	}

	if params.ExpiryTime != nil {
		now := time.Now()
		expMin := now.AddDate(0, 0, 28)
		expMax := now.AddDate(10, 0, 0)

		evaluatedExpTime := base.AddPeriod(now, *params.ExpiryTime)
		if evaluatedExpTime.Before(expMin) || evaluatedExpTime.After(expMax) {
			return fmt.Errorf("%w: expiry time cannot be more than 10 years or less than 28 days", base.ErrResponseStatusBadRequest)
		}
		doc.ExpiryTime = params.ExpiryTime
		digester.Write(doc.ExpiryTime.Bytes())
	}

	doc.Version = base.HexDigest(digester.Sum(nil))
	return nil
}

func (doc *KeyPolicyDoc) populateModelRef(m *KeyPolicyRef) {
	if doc == nil || m == nil {
		return
	}
	doc.BaseDoc.PopulateModelRef(&m.ResourceReference)
	m.DisplayName = doc.DisplayName
}

func (doc *KeyPolicyDoc) populateModel(m *KeyPolicy) {
	if doc == nil || m == nil {
		return
	}
	doc.populateModelRef(&m.KeyPolicyRef)
	m.KeyProperties = doc.KeyProperties
	m.Exportable = doc.Exportable
	m.ExpiryTime = doc.ExpiryTime
}
