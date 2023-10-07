package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertKeySpec struct {
	Alg     models.JwkAlg  `json:"alg"`
	Kty     models.JwtKty  `json:"kty"`
	KeySize *int           `json:"key_size,omitempty"`
	Crv     *models.JwtCrv `json:"crv,omitempty"`
}

// the Alg field of the second argument dk of default cert keyspec is the default algoritm for RSA only, as EC has fixed signing algorithm for the specified curve
func (k *CertKeySpec) initWithCreateTemplateInput(r *models.JwkProperties, dk CertKeySpec) {

	k.Kty = dk.Kty
	switch k.Kty {
	case models.KeyTypeRSA:
		k.KeySize = dk.KeySize
		k.Crv = nil
	case models.KeyTypeEC:
		k.Crv = dk.Crv
		k.KeySize = nil
	}

	if r == nil {
		return
	}

	switch r.Kty {
	case models.KeyTypeRSA:
		k.Kty = models.KeyTypeRSA
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
			case models.AlgRS256,
				models.AlgRS384,
				models.AlgRS512:
				// ok
			default:
				k.Alg = dk.Alg
			}
		}
	case models.KeyTypeEC:
		k.Kty = models.KeyTypeEC
		k.Crv = r.Crv
		if k.Crv == nil {
			k.Crv = dk.Crv
		} else {
			switch *k.Crv {
			case models.CurveNameP256,
				models.CurveNameP384:
				// ok
			default:
				k.Crv = dk.Crv
			}
		}
	}

	// decide alg for ec
	if k.Kty == models.KeyTypeEC {
		switch *k.Crv {
		case models.CurveNameP256:
			k.Alg = models.AlgES256
		case models.CurveNameP384:
			k.Alg = models.AlgES384
		}
	}
}

func (k *CertKeySpec) PopulateKeyProperties(r *models.JwkProperties) {
	if k == nil || r == nil {
		return
	}
	r.Kty = k.Kty
	r.KeySize = k.KeySize
	r.Crv = k.Crv
	r.Alg = utils.ToPtr(k.Alg)
}
