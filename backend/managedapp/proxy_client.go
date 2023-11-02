package managedapp

import (
	"fmt"

	"github.com/rs/zerolog/log"
	agentproxyclient "github.com/stephenzsy/small-kms/backend/agent/proxyclient"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func (s *server) getProxiedClient(c ctx.RequestContext,
	instanceID base.Identifier, accessToken string) (agentproxyclient.ClientWithResponsesInterface, error) {
	logger := log.Ctx(c).With().Str("proc", "getProxiedClient").Logger()
	if client, tokenExpired, err := s.proxyClientPool.GetCachedClient(accessToken); err != nil {
		logger.Debug().Err(err).Msg("failed to get cached client")
		return nil, err
	} else if tokenExpired {
		logger.Debug().Int("poolSize", s.proxyClientPool.Len()).Msg("cache hit: token expired")
		return nil, fmt.Errorf("%w: token expired", base.ErrResponseStatusBadRequest)
	} else if client != nil {
		logger.Debug().Int("poolSize", s.proxyClientPool.Len()).Msg("cache hit: client pool")
		return client, nil
	}
	logger.Debug().Int("poolSize", s.proxyClientPool.Len()).Msg("cache miss: client pool")

	// resolve instanceID to server endpoint
	doc, err := apiReadAgentInstanceDoc(c, instanceID)
	if err != nil {
		return nil, err
	}
	if doc.Endpoint == "" {
		return nil, fmt.Errorf("%w: agent instance has no endpoint", base.ErrResponseStatusBadRequest)
	}

	configDoc, err := apiReadAgentConfigDoc(c)
	if err != nil {
		return nil, err
	}

	certDoc, err := cert.ReadCertDocByID(c, configDoc.TLSCertificateID)
	if err != nil {
		return nil, err
	}
	certChain := certDoc.KeySpec.CertificateChain

	// create new client
	client, err := agentproxyclient.NewClientWithCreds(doc.Endpoint, certChain[len(certChain)-1], accessToken)
	if err != nil {
		return nil, err
	}
	if err := s.proxyClientPool.AddClient(client, accessToken); err != nil {
		log.Error().Err(err).Msg("failed to add client to pool")
	} else {
		logger.Debug().Int("poolSize", s.proxyClientPool.Len()).Msg("cache add: new client")
	}
	return client, nil
}
