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
	client     *azkeys.Client
	c          context.Context
	kid        azkeys.ID
	publicKey  crypto.PublicKey
	keyType    i.JsonWebKeyType
	jwsa       i.JsonWebSignatureAlgorithm
	formatX509 bool
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
func (ck *azCloudKey) Sign(rand io.Reader, digest []byte, _ crypto.SignerOpts) (signature []byte, err error) {
	if azSignAlg, ok := algMapping[ck.jwsa]; !ok {
		return nil, fmt.Errorf("%w: %s", i.ErrInvalidAlgorithm, ck.jwsa)
	} else if resp, err := ck.client.Sign(ck.c, ck.kid.Name(), ck.kid.Version(), azkeys.SignParameters{
		Algorithm: &azSignAlg,
		Value:     digest,
	}, nil); err != nil {
		return nil, err
	} else if ck.formatX509 {
		return toX509Signature(resp.Result, azSignAlg)
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

func NewAzCloudSignatureKeyWithKID(c context.Context, client *azkeys.Client, kid string,
	jwsa i.JsonWebSignatureAlgorithm, formatX509 bool, publicKey crypto.PublicKey) i.CloudSignatureKey {
	return &azCloudKey{
		client:     client,
		kid:        azkeys.ID(kid),
		c:          c,
		jwsa:       jwsa,
		formatX509: formatX509,
		publicKey:  publicKey,
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

func CreateCloudSignatureKey(c context.Context,
	client *azkeys.Client,
	name string,
	ckParams azkeys.CreateKeyParameters,
	jwsa i.JsonWebSignatureAlgorithm,
	formatX509 bool) (azkeys.CreateKeyResponse, i.CloudSignatureKey, error) {
	ckResp, err := client.CreateKey(c, name, ckParams, nil)
	if err != nil {
		return ckResp, nil, err
	}

	return ckResp, &azCloudKey{
		client:     client,
		kid:        *ckResp.Key.KID,
		c:          c,
		publicKey:  newSigningJWKFromKeyVaultKey(ckResp.Key).PublicKey(),
		jwsa:       jwsa,
		formatX509: formatX509,
	}, nil
}
