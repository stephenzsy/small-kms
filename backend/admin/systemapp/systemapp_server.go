package systemapp

import (
	appadmin "github.com/stephenzsy/small-kms/backend/admin/app"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
)

type SystemAppName string

// Defines values for SystemAppName.
const (
	SystemAppNameAPI     SystemAppName = "api"
	SystemAppNameBackend SystemAppName = "backend"
)

type SystemAppDoc = appadmin.AppDoc

func validateSystemAppName(name string) (SystemAppName, error) {
	typed := SystemAppName(name)
	switch typed {
	case SystemAppNameAPI, SystemAppNameBackend:
		return typed, nil
	default:
		return typed, base.ErrResponseStatusNotFound
	}
}

type SystemAppAdminServer struct {
	api.APIServer
}

func NewServer(apiServer api.APIServer) *SystemAppAdminServer {
	return &SystemAppAdminServer{
		APIServer: apiServer,
	}
}
