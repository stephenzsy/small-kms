package key

import (
	"github.com/stephenzsy/small-kms/backend/api"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
)

type KeyAdminServer struct {
	api.APIServer
}

func NewServer(apiServer api.APIServer) *KeyAdminServer {
	return &KeyAdminServer{
		APIServer: apiServer,
	}
}

type JsonWebKeyOperation = cloudkey.JsonWebKeyOperation
