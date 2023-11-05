package common

type commonConfig struct {
	serviceIdentity AzureIdentity
	envService      EnvService
	buildID         string
}

// BuildID implements CommonServer.
func (c commonConfig) BuildID() string {
	return c.buildID
}

func NewCommonConfig(envSvc EnvService, buildID string) (c commonConfig, err error) {
	c.envService = envSvc
	c.buildID = buildID
	c.serviceIdentity, err = NewAzureIdentityFromEnv(envSvc, IdentityEnvVarPrefixService)
	if err != nil {
		return c, err
	}
	return
}

type CommonServer interface {
	ServiceIdentityProvider
	EnvService() EnvService
	BuildID() string
}

// ServiceIdentity implements CommonServer.
func (c commonConfig) ServiceIdentity() AzureIdentity {
	return c.serviceIdentity
}

func (c commonConfig) EnvService() EnvService {
	return c.envService
}

var _ CommonServer = commonConfig{}
