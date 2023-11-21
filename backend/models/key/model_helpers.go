package keymodels

import (
	"io"
	"strconv"

	"github.com/stephenzsy/small-kms/backend/models"
)

type (
	keyPolicyComposed struct {
		models.Ref
		KeyPolicyFields
	}

	keyRefComposed struct {
		models.Ref
		KeyRefFields
	}

	keyComposed struct {
		KeyRef
		KeyFields
	}
)

func (jwkspec *JsonWebKeySpec) Digest(w io.Writer) {
	w.Write([]byte(jwkspec.Kty))
	if jwkspec.KeySize != nil {
		io.WriteString(w, strconv.Itoa(*jwkspec.KeySize))
	}
	w.Write([]byte(jwkspec.Crv))
	for _, op := range jwkspec.KeyOperations {
		w.Write([]byte(op))
	}
}
