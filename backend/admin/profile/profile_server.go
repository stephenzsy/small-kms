package profile

import "github.com/stephenzsy/small-kms/backend/api"

type ProfileServer struct {
	api.APIServer
}

func NewServer(apiServer api.APIServer) *ProfileServer {
	return &ProfileServer{
		APIServer: apiServer,
	}
}
