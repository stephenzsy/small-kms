package systemapp

import (
	"github.com/labstack/echo/v4"
	appadmin "github.com/stephenzsy/small-kms/backend/admin/app"
	"github.com/stephenzsy/small-kms/backend/api"
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
		return typed, echo.ErrNotFound
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
