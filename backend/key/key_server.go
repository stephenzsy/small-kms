package key

import (
	"github.com/stephenzsy/small-kms/backend/api"
)

type server struct {
	api.APIServer
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
