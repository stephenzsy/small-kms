package cert

import (
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
)

type CertServer struct {
	api.APIServer
	cryptoProvider cryptoprovider.CryptoProvider
}

func NewServer(apiServer api.APIServer) (*CertServer, error) {
	cryptoStore, err := cryptoprovider.NewCryptoProvider()
	if err != nil {
		return nil, err
	}

	return &CertServer{
		APIServer:      apiServer,
		cryptoProvider: cryptoStore,
	}, nil
}
