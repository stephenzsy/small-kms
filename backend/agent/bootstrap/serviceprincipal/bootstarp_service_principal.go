package serviceprincipal

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type ServicePrincipalBootstraper struct {
}

func NewServicePrincipalBootstraper() *ServicePrincipalBootstraper {
	return &ServicePrincipalBootstraper{}
}

func (*ServicePrincipalBootstraper) Bootstrap(c context.Context, certPath string) error {
	if certPath == "" {
		return errors.New("missing client cert path")
	}
	if _, err := os.Stat(certPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		fmt.Println("client cert already exists, skipping")
		return nil
	}
	return nil
}
