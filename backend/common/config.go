package common

import (
	"log"
	"os"
)

type ServerRole string

const ()

type commonConfig struct {
	serviceIdentity AzureIdentity
}

func NewCommonConfig() (c commonConfig, err error) {
	c.serviceIdentity, err = NewAzureIdentityFromEnv(IdentityEnvVarPrefixService)
	if err != nil {
		return
	}

	return
}

type CommonServer interface {
	ServiceIdentityProvider
}

// ServiceIdentity implements CommonServer.
func (c commonConfig) ServiceIdentity() AzureIdentity {
	return c.serviceIdentity
}

var _ CommonServer = commonConfig{}

func MustGetenv(name string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		log.Panicf("No variable %s configured", name)
	}
	log.Printf("Config %s = %s", name, value)
	return
}

func MustGetenvSecret(name string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		log.Panicf("No variable %s configured", name)
	}
	log.Printf("Config %s = **********", name)
	return
}

func GetEnvWithDefault(name string, defaultValue string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		value = defaultValue
	}
	log.Printf("Config %s = %s", name, value)
	return
}
