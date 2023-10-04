package secret

type WindowsSecretsService struct{}

// GenerateRSAKey implements SecretsService.
func (*WindowsSecretsService) GenerateRSAKey() {
	cng.GenerateRSAKey()
	panic("unimplemented")
}
