package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertKeySpec struct {
	Alg     shared.JwkAlg  `json:"alg"`
	Kty     shared.JwtKty  `json:"kty"`
	KeySize *int32         `json:"key_size,omitempty"`
	Crv     *shared.JwkCrv `json:"crv,omitempty"`
}

// the Alg field of the second argument dk of default cert keyspec is the default algoritm for RSA only, as EC has fixed signing algorithm for the specified curve
func (k *CertKeySpec) initWithCreateTemplateInput(r *shared.JwkProperties, dk CertKeySpec) {

	k.Kty = dk.Kty
	switch k.Kty {
	case shared.KeyTypeRSA:
		k.KeySize = dk.KeySize
		k.Crv = nil
		k.Alg = dk.Alg
	case shared.KeyTypeEC:
		k.Crv = dk.Crv
		k.KeySize = nil
	}

	if r != nil {
		switch r.Kty {
		case shared.KeyTypeRSA:
			k.Kty = shared.KeyTypeRSA
			k.KeySize = r.KeySize
			k.Alg = dk.Alg
			if k.KeySize == nil {
				k.KeySize = dk.KeySize
				switch *k.KeySize {
				case 2048, 3072, 4096:
					// ok
				default:
					k.KeySize = dk.KeySize
				}
			}
			if r.Alg != nil {
				k.Alg = *r.Alg
				switch k.Alg {
				case shared.AlgRS256,
					shared.AlgRS384,
					shared.AlgRS512:
					// ok
				default:
					k.Alg = dk.Alg
				}
			}
		case shared.KeyTypeEC:
			k.Kty = shared.KeyTypeEC
			k.Crv = r.Crv
			if k.Crv == nil {
				k.Crv = dk.Crv
			} else {
				switch *k.Crv {
				case shared.CurveNameP256,
					shared.CurveNameP384:
					// ok
				default:
					k.Crv = dk.Crv
				}
			}
		}
	}

	// decide alg for ec
	if k.Kty == shared.KeyTypeEC {
		switch *k.Crv {
		case shared.CurveNameP256:
			k.Alg = shared.AlgES256
		case shared.CurveNameP384:
			k.Alg = shared.AlgES384
		}
	}
}

func (k *CertKeySpec) PopulateKeyProperties(r *shared.JwkProperties) {
	if k == nil || r == nil {
		return
	}
	r.Kty = k.Kty
	r.KeySize = k.KeySize
	r.Crv = k.Crv
	r.Alg = utils.ToPtr(k.Alg)
}
