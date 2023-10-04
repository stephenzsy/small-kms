package secret

import "runtime"

type SecretsService interface {
	GenerateRSAKey()
}

func GetService() SecretsService {
	if runtime.GOOS == "windows" {
		return &WindowsSecretsService{}
	}
	return nil
}
