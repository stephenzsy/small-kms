package secret

import (
	ctx "context"
	"crypto/rsa"
	"runtime"
)

type keySessionContextKey string

type SecretService interface {
	RS256SignHash(hash []byte, keyIdentifier string) (signature []byte, publicKey *rsa.PublicKey, err error)
}

func GetService(context ctx.Context) SecretService {
	if runtime.GOOS == "windows" {
		return &WindowsSecretsService{}
	}
	return nil
}
