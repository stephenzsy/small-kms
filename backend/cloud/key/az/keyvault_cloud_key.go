package cloudkeyaz

import (
	"context"
	"crypto"
	"crypto/rsa"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	i "github.com/stephenzsy/small-kms/backend/cloud/key"
)

type azCloudKey struct {
	client    *azkeys.Client
	c         context.Context
	kid       azkeys.ID
	publicKey crypto.PublicKey
	keyType   i.JsonWebKeyType
}

var keyOpMapping = map[i.JsonWebKeyOperation]azkeys.KeyOperation{
	i.JsonWebKeyOperationSign:      azkeys.KeyOperationSign,
	i.JsonWebKeyOperationVerify:    azkeys.KeyOperationVerify,
	i.JsonWebKeyOperationEncrypt:   azkeys.KeyOperationEncrypt,
	i.JsonWebKeyOperationDecrypt:   azkeys.KeyOperationDecrypt,
	i.JsonWebKeyOperationWrapKey:   azkeys.KeyOperationWrapKey,
	i.JsonWebKeyOperationUnwrapKey: azkeys.KeyOperationUnwrapKey,
}

var algMapping = map[i.JsonWebSignatureAlgorithm]azkeys.SignatureAlgorithm{
	i.SignatureAlgorithmRS256: azkeys.SignatureAlgorithmRS256,
	i.SignatureAlgorithmRS384: azkeys.SignatureAlgorithmRS384,
	i.SignatureAlgorithmRS512: azkeys.SignatureAlgorithmRS512,
	i.SignatureAlgorithmPS256: azkeys.SignatureAlgorithmPS256,
	i.SignatureAlgorithmPS384: azkeys.SignatureAlgorithmPS384,
	i.SignatureAlgorithmPS512: azkeys.SignatureAlgorithmPS512,
	i.SignatureAlgorithmES256: azkeys.SignatureAlgorithmES256,
	i.SignatureAlgorithmES384: azkeys.SignatureAlgorithmES384,
	i.SignatureAlgorithmES512: azkeys.SignatureAlgorithmES512,
}

var ktyReverseMapping = map[azkeys.KeyType]i.JsonWebKeyType{
	azkeys.KeyTypeEC:  i.KeyTypeEC,
	azkeys.KeyTypeRSA: i.KeyTypeRSA,
	azkeys.KeyTypeOct: i.KeyTypeOct,
}

var crvReverseMapping = map[azkeys.CurveName]i.JsonWebKeyCurveName{
	azkeys.CurveNameP256: i.CurveNameP256,
	azkeys.CurveNameP384: i.CurveNameP384,
	azkeys.CurveNameP521: i.CurveNameP521,
}

func (ck *azCloudKey) KeyType() i.JsonWebKeyType {
	return ck.keyType
}

func (ck *azCloudKey) KeyID() string {
	return string(ck.kid)
}

func newSigningJWKFromKeyVaultKey(kvKey *azkeys.JSONWebKey) *i.JsonWebKey[i.JsonWebSignatureAlgorithm] {
	keyID := string(*kvKey.KID)
	r := &i.JsonWebKey[i.JsonWebSignatureAlgorithm]{
		JsonWebKeyBase: i.JsonWebKeyBase{
			KeyType: ktyReverseMapping[*kvKey.Kty],
			KeyID:   keyID,
			X:       kvKey.X,
			Y:       kvKey.Y,
			N:       kvKey.N,
			E:       kvKey.E,
		},
	}
	if kvKey.Crv != nil {
		crv := crvReverseMapping[*kvKey.Crv]
		r.Curve = crv
	}
	return r
}

func ToAzKeyOperation(op i.JsonWebKeyOperation) azkeys.KeyOperation {
	if v, ok := keyOpMapping[op]; !ok {
		return azkeys.KeyOperation("")
	} else {
		return v
	}
}

// Public implements cloudkey.CloudSignatureKey.
func (ck *azCloudKey) Public() crypto.PublicKey {
	if ck.publicKey == nil {
		resp, err := ck.client.GetKey(ck.c, ck.kid.Name(), ck.kid.Version(), nil)
		if err != nil {
			ck.publicKey = err
		}
		ck.publicKey = newSigningJWKFromKeyVaultKey(resp.Key).PublicKey()
	}
	return ck.publicKey
}

// Sign implements cloudkey.CloudSignatureKey.
func (ck *azCloudKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if alg, ok := opts.(i.JsonWebSignatureAlgorithm); !ok {
		return nil, fmt.Errorf("%w: %T", i.ErrInvalidAlgorithm, opts)
	} else if azSignAlg, ok := algMapping[alg]; !ok {
		return nil, fmt.Errorf("%w: %s", i.ErrInvalidAlgorithm, alg)
	} else if resp, err := ck.client.Sign(ck.c, ck.kid.Name(), ck.kid.Version(), azkeys.SignParameters{
		Algorithm: &azSignAlg,
		Value:     digest,
	}, nil); err != nil {
		return nil, err
	} else {
		return resp.Result, nil
	}
}

// Use unwrap key, large blob of data should not be decrypted directly with this key
func (ck *azCloudKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	params := azkeys.KeyOperationParameters{
		Value: msg,
	}
	switch ck.keyType {
	case i.KeyTypeRSA:
		if opts, ok := opts.(*rsa.OAEPOptions); !ok {
			return nil, fmt.Errorf("%w: %T", i.ErrInvalidAlgorithm, opts)
		} else {
			switch opts.Hash {
			case crypto.SHA256:
				params.Algorithm = to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP256)
			case crypto.SHA1:
				params.Algorithm = to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP)
			default:
				return nil, fmt.Errorf("%w: %s", i.ErrInvalidAlgorithm, opts.Hash)
			}
		}
	default:
		return nil, fmt.Errorf("%w: %s", i.ErrInvalidKeyType, ck.keyType)
	}
	resp, err := ck.client.UnwrapKey(ck.c, ck.kid.Name(), ck.kid.Version(), params, nil)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

var _ i.CloudSignatureKey = (*azCloudKey)(nil)
var _ i.CloudWrappingKey = (*azCloudKey)(nil)

func NewAzCloudSignatureKeyWithKID(c context.Context, client *azkeys.Client, kid string) i.CloudSignatureKey {
	return &azCloudKey{
		client: client,
		kid:    azkeys.ID(kid),
		c:      c,
	}
}

func NewCloudWrappingKeyWithKID(c context.Context, client *azkeys.Client, kid string, keyType i.JsonWebKeyType) i.CloudWrappingKey {
	return &azCloudKey{
		client:  client,
		kid:     azkeys.ID(kid),
		c:       c,
		keyType: keyType,
	}
}
