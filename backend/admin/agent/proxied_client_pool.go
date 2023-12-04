package agentadmin

import (
	"container/heap"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	nethttp "net/http"
	"sync"
	"time"

	agentauth "github.com/stephenzsy/small-kms/backend/agent/auth"
	agentendpoint "github.com/stephenzsy/small-kms/backend/agent/endpoint"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	cloudkeyx "github.com/stephenzsy/small-kms/backend/cloud/key/x"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
)

type pooledClient struct {
	agentendpoint.ClientWithResponsesInterface
	exp         time.Time
	namespaceID string
	instanceID  string
}

type ProxyClientPool struct {
	items   []*pooledClient
	lock    sync.RWMutex
	maxSize int
	cutoff  time.Time
}

func newClientWithCreds(server string, caCertBytes []byte, accessToken string) (agentendpoint.ClientWithResponsesInterface, error) {
	caPool := x509.NewCertPool()
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		return nil, err
	}
	caPool.AddCert(caCert)
	clientTransport := nethttp.DefaultTransport.(*nethttp.Transport).Clone()
	clientTransport.ForceAttemptHTTP2 = true
	clientTransport.DisableCompression = false
	clientTransport.TLSClientConfig = &tls.Config{
		RootCAs:    caPool,
		MinVersion: tls.VersionTLS13,
	}
	return agentendpoint.NewClientWithResponses(server, agentendpoint.WithHTTPClient(&nethttp.Client{
		Transport: clientTransport,
		CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
			return nethttp.ErrUseLastResponse
		},
	}), agentendpoint.WithRequestEditorFn(func(ctx context.Context, req *nethttp.Request) error {
		if req.Header.Get("Authorization") == "" {
			req.Header.Set("Authorization", "Bearer "+accessToken)
		}
		return nil
	}))
}

// Len implements heap.Interface.
func (p *ProxyClientPool) Len() int {
	return len(p.items)
}

// Less implements heap.Interface.
func (p *ProxyClientPool) Less(i int, j int) bool {
	return p.items[i].exp.Before(p.items[j].exp)
}

// Pop implements heap.Interface.
func (p *ProxyClientPool) Pop() any {
	item := p.items[len(p.items)-1]
	p.items[len(p.items)-1] = nil
	p.items = p.items[0 : len(p.items)-1]
	return item
}

// Push implements heap.Interface.
func (p *ProxyClientPool) Push(x any) {
	item := x.(*pooledClient)
	p.items = append(p.items, item)
}

// Swap implements heap.Interface.
func (p *ProxyClientPool) Swap(i int, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

var _ heap.Interface = (*ProxyClientPool)(nil)

func (p *ProxyClientPool) cleanup() {
	cutoff := time.Now()
	for len(p.items) > 0 && p.items[0].exp.Before(cutoff) {
		heap.Pop(p)
	}
}

func (p *ProxyClientPool) GetClient(c context.Context, namespaceID string, instanceID string) (agentendpoint.ClientWithResponsesInterface, error) {
	p.lock.Lock()
	p.cleanup()
	p.lock.Unlock()

	// claims := jwt.RegisteredClaims{}
	// _, _, err := tokenParser.ParseUnverified(tokenString, &claims)
	// if err != nil {
	// 	return nil, false, err
	// }
	// if claims.ExpiresAt.Before(p.GetCutOff()) {
	// 	return nil, true, nil
	// }
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, item := range p.items {
		if item.namespaceID == namespaceID && item.instanceID == instanceID {
			return item, nil
		}
	}
	if len(p.items) >= p.maxSize {
		return nil, base.ErrResponseStatusBadRequest
	}

	instanceDoc, err := getAgentInstanceInternal(c, namespaceID, instanceID)
	if err != nil {
		return nil, err
	} else if instanceDoc.TlsCertificateID == "" {
		return nil, errors.New("instance does not have tls certificate")
	} else if instanceDoc.JwtVerfyKeyID == "" {
		return nil, errors.New("instance does not have jwt verify key")
	} else if instanceDoc.Endpoint == "" {
		return nil, errors.New("instance does not have endpoint")
	}

	certDoc, err := cert.GetCertificateInternal(c, models.NamespaceProviderServicePrincipal, namespaceID, instanceDoc.TlsCertificateID)
	if err != nil {
		return nil, err
	}

	keyDoc, err := key.GetKeyInternal(c, models.NamespaceProviderServicePrincipal, namespaceID, instanceDoc.JwtVerfyKeyID)
	if err != nil {
		return nil, err
	}

	ck := cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, kv.GetAzKeyVaultService(c).AzKeysClient(), keyDoc.JsonWebKey.KeyID, cloudkey.SignatureAlgorithmES384, false, keyDoc.PublicKey())
	identity := auth.GetAuthIdentity(c)
	accessToken, exp, err := agentauth.NewSignedAgentAuthJWT(cloudkeyx.NewJWTSigningMethod(cloudkey.SignatureAlgorithmES384), identity.ClientPrincipalID().String(), instanceDoc.Endpoint, ck)
	if err != nil {
		return nil, err
	}

	certChain := certDoc.GetJsonWebKey().CertificateChain
	client, err := newClientWithCreds(instanceDoc.Endpoint, certChain[len(certChain)-1], accessToken)
	if err != nil {
		return nil, err
	}
	heap.Push(p, &pooledClient{
		ClientWithResponsesInterface: client,
		exp:                          exp,
		namespaceID:                  namespaceID,
		instanceID:                   instanceID,
	})
	return client, nil
}

// func (p *ProxyClientPool) AddClient(client ClientWithResponsesInterface, tokenString string) error {
// 	p.cleanup()
// 	claims := jwt.RegisteredClaims{}
// 	_, _, err := tokenParser.ParseUnverified(tokenString, &claims)
// 	if err != nil {
// 		return err
// 	}
// 	if claims.ExpiresAt.Before(p.GetCutOff()) {
// 		return errors.New("token is expired or pool is full for old tokens")
// 	}
// 	p.lock.Lock()
// 	defer p.lock.Unlock()
// 	if len(p.items) >= p.maxSize {
// 		if claims.ExpiresAt.Before(p.items[0].exp) || claims.ExpiresAt.Equal(p.items[0].exp) {
// 			p.cutoff = p.items[0].exp.Add(time.Nanosecond)
// 			// reject since client is too old
// 			return errors.New("pool is full, client token is too old")
// 		}
// 		// replace the oldest client
// 		p.items[0].client = client
// 		p.items[0].tokenID = claims.ID
// 		p.items[0].exp = claims.ExpiresAt.Time
// 		heap.Fix(p, 0)
// 	} else {
// 		heap.Push(p, &pooledClientItem{
// 			client:  client,
// 			tokenID: claims.ID,
// 			exp:     claims.ExpiresAt.Time,
// 		})
// 	}
// 	return nil
// }

func (p *ProxyClientPool) GetCutOff() time.Time {
	now := time.Now()
	if p.cutoff.Before(now) {
		p.cutoff = now
	}
	return p.cutoff
}

func NewProxyClientPool(maxSize int) *ProxyClientPool {
	return &ProxyClientPool{
		items:   make([]*pooledClient, 0, maxSize),
		maxSize: maxSize,
	}
}
