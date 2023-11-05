package common

type commonConfig struct {
	serviceIdentity AzureIdentity
	envService      EnvService
}

func NewCommonConfig(envSvc EnvService) (c commonConfig, err error) {
	c.envService = envSvc
	c.serviceIdentity, err = NewAzureIdentityFromEnv(envSvc, IdentityEnvVarPrefixService)
	if err != nil {
		return
	}

	return
}

type CommonServer interface {
	ServiceIdentityProvider
	EnvService() EnvService
}

// ServiceIdentity implements CommonServer.
func (c commonConfig) ServiceIdentity() AzureIdentity {
	return c.serviceIdentity
}

func (c commonConfig) EnvService() EnvService {
	return c.envService
}

var _ CommonServer = commonConfig{}
