package cert

import (
	"github.com/stephenzsy/small-kms/backend/api"
)

type CertServer struct {
	api.APIServer
}

func NewServer(apiServer api.APIServer) *CertServer {
	return &CertServer{
		APIServer: apiServer,
	}
}
