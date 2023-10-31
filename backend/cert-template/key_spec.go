package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/shared"
)

type CertKeySpec struct {
	Alg     shared.JwkAlg  `json:"alg"`
	Kty     shared.JwtKty  `json:"kty"`
	KeySize *int32         `json:"key_size,omitempty"`
	Crv     *shared.JwkCrv `json:"crv,omitempty"`
}

func (k *CertKeySpec) PopulateKeyProperties(r *shared.JwkProperties) {
	if k == nil || r == nil {
		return
	}
	r.Alg = &k.Alg
	r.Kty = k.Kty
	r.KeySize = k.KeySize
	r.Crv = k.Crv
}
