package key

import (
	"io"
	"strconv"
)

func (p *GenerateJsonWebKeyProperties) writeToDigest(w io.Writer) {
	if p == nil {
		return
	}
	w.Write([]byte(p.Kty))
	if p.KeySize != nil {
		w.Write([]byte(strconv.Itoa(int(*p.KeySize))))
	}
	w.Write([]byte(p.Crv))
	for _, op := range p.KeyOperations {
		w.Write([]byte(op))
	}
}
