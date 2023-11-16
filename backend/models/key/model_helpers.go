package keymodels

import (
	"io"
	"strconv"

	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type (
	keyPolicyComposed struct {
		models.Ref
		KeyPolicyFields
	}
)

func (jwkspec *JsonWebKeySpec) Digest(w io.Writer) (int, error) {
	cw := utils.ChainedWriter{InnerWriter: w}
	cw.Write([]byte(jwkspec.Kty))
	if jwkspec.KeySize != nil {
		cw.WriteString(strconv.Itoa(*jwkspec.KeySize))
	}
	cw.Write([]byte(jwkspec.Crv))
	for _, op := range jwkspec.KeyOperations {
		cw.Write([]byte(op))
	}
	return cw.Return()
}
