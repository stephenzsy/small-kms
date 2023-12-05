package key

import (
	"github.com/stephenzsy/small-kms/backend/api"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
)

type KeyAdminServer struct {
	api.APIServer
	cryptoProvider cryptoprovider.CryptoProvider
}

func NewServer(apiServer api.APIServer) (*KeyAdminServer, error) {
	cryptoProvider, err := cryptoprovider.NewCryptoProvider()
	if err != nil {
		return nil, err
	}
	return &KeyAdminServer{
		APIServer:      apiServer,
		cryptoProvider: cryptoProvider,
	}, nil
}

type JsonWebKeyOperation = cloudkey.JsonWebKeyOperation
