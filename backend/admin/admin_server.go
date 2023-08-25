package admin

import (
	"github.com/stephenzsy/small-kms/backend/common"
)

type adminServer struct {
	config common.ServerConfig
}

func NewAdminServer(c common.ServerConfig) ServerInterface {
	return &adminServer{config: c}
}
