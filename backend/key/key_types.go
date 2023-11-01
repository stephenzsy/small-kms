package key

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type (
	keyPolicyRefComposed struct {
		base.ResourceReference
		KeyPolicyRefFields
	}

	keyPolicyComposed struct {
		KeyPolicyRef
		KeyPolicyFields
	}

	keyComposed struct {
		base.ResourceReference
		KeySpec
		KeyFields
	}
)

func (ks *SigningKeySpec) WriteToDigest(w io.Writer) (s int, err error) {
	if ks == nil {
		return 0, nil
	}
	if ks.Alg != nil {
		if c, err := w.Write([]byte(*ks.Alg)); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	if c, err := w.Write([]byte(ks.Kty)); err != nil {
		return s + c, err
	} else {
		s += c
	}
	switch ks.Kty {
	case "RSA":
		if ks.KeySize != nil {
			if c, err := w.Write([]byte(strconv.Itoa(int(*ks.KeySize)))); err != nil {
				return s + c, err
			} else {
				s += c
			}
		}
	case "EC":
		if ks.Crv != nil {
			if c, err := w.Write([]byte(*ks.Crv)); err != nil {
				return s + c, err
			} else {
				s += c
			}
		}
	}
	for _, op := range ks.KeyOperations {
		if c, err := w.Write([]byte(op)); err != nil {
			return s + c, err
		}
	}
	return s, nil
}

func (la *LifetimeAction) WriteToDigest(w io.Writer) (s int, err error) {
	if la == nil {
		return 0, nil
	}
	return la.Trigger.WriteToDigest(w)
}

func (lt *LifetimeTrigger) WriteToDigest(w io.Writer) (s int, err error) {
	if lt == nil {
		return 0, nil
	}
	if lt.TimeAfterCreate != nil {
		if c, err := w.Write([]byte("timeAfterCreate")); err != nil {
			return s + c, err
		} else {
			s += c
		}
		if c, err := w.Write(lt.TimeAfterCreate.Bytes()); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	if lt.TimeBeforeExpiry != nil {
		if c, err := w.Write([]byte("timeBeforeExpiry")); err != nil {
			return s + c, err
		} else {
			s += c
		}
		if c, err := w.Write(lt.TimeBeforeExpiry.Bytes()); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	if lt.PercentageAfterCreate != nil {
		if c, err := w.Write([]byte("percentageAfterCreate")); err != nil {
			return s + c, err
		} else {
			s += c
		}
		if c, err := w.Write([]byte(strconv.Itoa(int(*lt.PercentageAfterCreate)))); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	return s, nil
}

func (ks *SigningKeySpec) SetPublicKey(pubKey crypto.PublicKey) error {
	switch pubKey := pubKey.(type) {
	case *rsa.PublicKey:
		ks.Kty = cloudkey.KeyTypeRSA
		ks.N = base.Base64RawURLEncodedBytes(pubKey.N.Bytes())
		ks.E = base.Base64RawURLEncodedBytes(big.NewInt(int64(pubKey.E)).Bytes())
	case *ecdsa.PublicKey:
		ks.Kty = cloudkey.KeyTypeEC
		ks.X = base.Base64RawURLEncodedBytes(pubKey.X.Bytes())
		ks.Y = base.Base64RawURLEncodedBytes(pubKey.Y.Bytes())
		switch pubKey.Curve.Params().Name {
		case elliptic.P256().Params().Name:
			ks.Crv = utils.ToPtr(cloudkey.CurveNameP256)
		case elliptic.P384().Params().Name:
			ks.Crv = utils.ToPtr(cloudkey.CurveNameP384)
		case elliptic.P521().Params().Name:
			ks.Crv = utils.ToPtr(cloudkey.CurveNameP521)
		default:
			return fmt.Errorf("unsupported EC curve: %s", pubKey.Curve.Params().Name)
		}
	default:
		return fmt.Errorf("unsupported public key type: %T", pubKey)
	}
	return nil
}
