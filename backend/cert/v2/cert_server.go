package cert

import (
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
)

type CertServer struct {
	api.APIServer
	cryptoStore cryptoprovider.CryptoProvider
}

func NewServer(apiServer api.APIServer) *CertServer {
	cryptoStore, _ := cryptoprovider.NewCryptoProvider()

	return &CertServer{
		APIServer:   apiServer,
		cryptoStore: cryptoStore,
	}
}
