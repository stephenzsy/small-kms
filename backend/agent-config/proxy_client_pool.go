package agentconfig

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/admin/agentproxyclient"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type AgentProxyHttpClientPool interface {
	GetProxyHttpClient(eCtx context.Context, nsID shared.NamespaceIdentifier) (agentproxyclient.ClientWithResponsesInterface, error)
}

type agentProxyHttpClientPool struct {
	cachedClients map[uuid.UUID]agentproxyclient.ClientWithResponsesInterface
	cacheLock     sync.RWMutex
}

type wrappedHttpClient struct {
	http.Client
	cleanupTimer *time.Timer
}

func (c *wrappedHttpClient) Do(req *http.Request) (*http.Response, error) {
	c.cleanupTimer.Stop()
	defer c.cleanupTimer.Reset(5 * time.Minute)
	return c.Client.Do(req)
}

var _ agentproxyclient.HttpRequestDoer = (*wrappedHttpClient)(nil)

// GetProxyHttpClient implements AgentProxyHttpClientPool.
func (p *agentProxyHttpClientPool) GetProxyHttpClient(eCtx context.Context, nsID shared.NamespaceIdentifier) (agentproxyclient.ClientWithResponsesInterface, error) {
	nsUUID := nsID.Identifier().UUID()
	if nsUUID == uuid.Nil {
		return nil, common.ErrStatusBadRequest
	}

	p.cacheLock.RLock()
	client, ok := p.cachedClients[nsUUID]
	p.cacheLock.RUnlock()
	if ok {
		return client, nil
	}

	// load doc
	d := AgentActiveServerDoc{}
	err := kmsdoc.Read(eCtx, NewConfigDocLocator(nsID, shared.AgentConfigNameActiveServer), &d)
	if err != nil {
		return nil, err
	}
	clientCerts := make([]tls.Certificate, 0, 2)
	for i, certId := range d.AuthorizedCertificateIDs {
		if i > 1 {
			break
		}
		certDoc, err := cert.ReadCertDocByLocator(eCtx, shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCert, certId)))
		if err != nil {
			return nil, err
		}
		clientCert, err := cert.CertDocKeyPair(eCtx, certDoc)
		if err != nil {
			return nil, err
		}
		clientCerts = append(clientCerts, clientCert)
	}
	caCertPool := x509.NewCertPool()
	for _, c := range clientCerts {
		for _, cBytes := range c.Certificate[1:] {
			caCert, err := x509.ParseCertificate(cBytes)
			if err != nil {
				return nil, err
			}
			caCertPool.AddCert(caCert)
		}
	}

	httpClient := &wrappedHttpClient{
		Client: http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost:     2, // limit number of connection to reduce remote handshake signing
				TLSHandshakeTimeout: 30 * time.Second,
				IdleConnTimeout:     5 * time.Minute,
				TLSClientConfig: &tls.Config{
					Certificates: clientCerts,
					RootCAs:      caCertPool,
				},
			},
		},
		cleanupTimer: time.NewTimer(5 * time.Minute),
	}
	client, err = agentproxyclient.NewClientWithResponses(d.EndpointURL, agentproxyclient.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	go func() {
		<-httpClient.cleanupTimer.C
		p.cacheLock.Lock()
		delete(p.cachedClients, nsUUID)
		p.cacheLock.Unlock()
	}()
	p.cacheLock.Lock()
	defer p.cacheLock.Unlock()
	p.cachedClients[nsUUID] = client
	return client, nil
}

var _ AgentProxyHttpClientPool = (*agentProxyHttpClientPool)(nil)

func NewAgentProxyHttpClientPool() AgentProxyHttpClientPool {
	return &agentProxyHttpClientPool{
		cachedClients: make(map[uuid.UUID]agentproxyclient.ClientWithResponsesInterface),
	}
}

type contextKey string

const proxyHttpClientPoolContextKey contextKey = "proxyHttpClientPool"

func WithNewProxyHttpCLientPool(c context.Context) context.Context {
	return context.WithValue(c, proxyHttpClientPoolContextKey, NewAgentProxyHttpClientPool())
}
