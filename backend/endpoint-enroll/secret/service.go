package secret

import (
	ctx "context"
	"crypto/rsa"
)

type keySessionContextKey string

type SecretService interface {
	RS256SignHash(hash []byte, keyIdentifier string, installToMachine bool) (signature []byte, publicKey *rsa.PublicKey, err error)
}

func GetService(context ctx.Context) SecretService { /*
		if runtime.GOOS == "windows" {
			return &WindowsSecretsService{}
		}*/
	return nil
}
