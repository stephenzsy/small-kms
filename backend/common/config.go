package common

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type ServerRole string

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

func GetNonEmptyEnv(name string) (string, error) {
	value := os.Getenv(name)
	value = strings.TrimSpace(value)
	if value == "" {
		return value, fmt.Errorf("%w:%s", ErrMissingEnvVar, name)
	}
	log.Printf("Config %s = %s", name, value)
	return value, nil
}

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
