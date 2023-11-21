package key

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/models"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

type KeyPolicyDoc struct {
	resdoc.ResourceDoc
	DisplayName string `json:"displayName"`

	KeySpec    keymodels.JsonWebKeySpec `json:"keySpec"`
	Exportable bool                     `json:"exportable"`
	ExpiryTime *caldur.CalendarDuration `json:"expiryTime,omitempty"`

	Version []byte `json:"version"`
}

const (
	queryColumnDisplayName = "c.displayName"
)

func sanitizeKeyOperations(keyOps []JsonWebKeyOperation) []JsonWebKeyOperation {
	if len(keyOps) == 0 {
		return nil
	}

	// remove duplicates
	seen := make(map[JsonWebKeyOperation]bool)
	for _, keyOp := range keyOps {
		switch keyOp {
		case cloudkey.JsonWebKeyOperationSign, cloudkey.JsonWebKeyOperationVerify,
			cloudkey.JsonWebKeyOperationEncrypt, cloudkey.JsonWebKeyOperationDecrypt,
			cloudkey.JsonWebKeyOperationWrapKey, cloudkey.JsonWebKeyOperationUnwrapKey:
			seen[keyOp] = true
		}
	}

	result := make([]JsonWebKeyOperation, 0, len(seen))
	for keyOp := range seen {
		result = append(result, keyOp)
	}
	return result
}

func (doc *KeyPolicyDoc) init(c context.Context, req *keymodels.CreateKeyPolicyRequest) error {
	logger := log.Ctx(c)

	if req.DisplayName != "" {
		doc.DisplayName = req.DisplayName
	} else {
		doc.DisplayName = doc.ID
	}

	digester := md5.New()
	// key properties default
	doc.KeySpec = keymodels.JsonWebKeySpec{
		Kty:     cloudkey.KeyTypeRSA,
		KeySize: utils.ToPtr(2048),
		KeyOperations: []cloudkey.JsonWebKeyOperation{
			cloudkey.JsonWebKeyOperationSign,
			cloudkey.JsonWebKeyOperationVerify,
		},
	}
	if req.KeySpec != nil {
		switch req.KeySpec.Kty {
		case cloudkey.KeyTypeRSA:
			doc.KeySpec.Kty = cloudkey.KeyTypeRSA
			doc.KeySpec.KeySize = utils.ToPtr(2048)
			if req.KeySpec.KeySize != nil {
				switch *req.KeySpec.KeySize {
				case 2048, 3072, 4096:
					doc.KeySpec.KeySize = req.KeySpec.KeySize
				case 0:
					// use default
				default:
					logger.Warn().Int("keySize", *req.KeySpec.KeySize).Msg("invalid key size, default to 2048")
				}
			}

		case cloudkey.KeyTypeEC:
			doc.KeySpec.Kty = cloudkey.KeyTypeEC
			doc.KeySpec.KeySize = nil
			doc.KeySpec.Crv = cloudkey.CurveNameP384 // default to P384 as cost in KeyVault is the same as 256

			switch req.KeySpec.Crv {
			case cloudkey.CurveNameP256, cloudkey.CurveNameP384, cloudkey.CurveNameP521:
				doc.KeySpec.Crv = req.KeySpec.Crv
			case "":
				// default
			default:
				logger.Warn().Str("crv", string(req.KeySpec.Crv)).Msg("invalid curve, default to P384")
			}
		}

		keyOps := sanitizeKeyOperations(req.KeySpec.KeyOperations)
		if len(keyOps) > 0 {
			doc.KeySpec.KeyOperations = keyOps
		}
	}
	doc.KeySpec.Digest(digester)

	if req.Exportable != nil && *req.Exportable {
		doc.Exportable = true
		digester.Write([]byte("exportable"))
	}

	if req.ExpiryTime != "" {
		expTime, err := caldur.Parse(req.ExpiryTime)
		if err != nil {
			return fmt.Errorf("%w: invalid expiry time format", base.ErrResponseStatusBadRequest)
		}

		now := time.Now()
		expMin := now.AddDate(0, 0, 28)
		expMax := now.AddDate(10, 0, 0)

		evaluatedExpTime := caldur.Shift(now, expTime)
		if evaluatedExpTime.Before(expMin) || evaluatedExpTime.After(expMax) {
			return fmt.Errorf("%w: expiry time cannot be more than 10 years or less than 28 days", base.ErrResponseStatusBadRequest)
		}
		doc.ExpiryTime = &expTime
		digester.Write(doc.ExpiryTime.Bytes())
	}

	doc.Version = digester.Sum(nil)
	return nil
}

func (doc *KeyPolicyDoc) ToRef() (m models.Ref) {
	m = doc.ResourceDoc.ToRef()
	m.DisplayName = &doc.DisplayName
	return m
}

func (doc *KeyPolicyDoc) ToModel() (m keymodels.KeyPolicy) {
	m.Ref = doc.ResourceDoc.ToRef()
	m.KeySpec = doc.KeySpec
	m.Exportable = doc.Exportable
	if doc.ExpiryTime != nil {
		m.ExpiryTime = doc.ExpiryTime.String()
	}
	return m
}
