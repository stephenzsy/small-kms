package agentpush

import (
	"container/heap"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	nethttp "net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewClientWithCreds(server string, caCertBytes []byte, accessToken string) (ClientWithResponsesInterface, error) {
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
	return NewClientWithResponses(server, WithHTTPClient(&nethttp.Client{
		Transport: clientTransport,
		CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
			return nethttp.ErrUseLastResponse
		},
		Timeout: time.Second * 100,
	}), WithRequestEditorFn(func(ctx context.Context, req *nethttp.Request) error {
		if req.Header.Get("Authorization") == "" {
			req.Header.Set("Authorization", "Bearer "+accessToken)
		}
		return nil
	}))
}

type pooledClientItem struct {
	client  ClientWithResponsesInterface
	tokenID string
	exp     time.Time
}

type ProxyClientPool struct {
	items   []*pooledClientItem
	lock    sync.RWMutex
	maxSize int
	cutoff  time.Time
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
	item := x.(*pooledClientItem)
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
		p.lock.Lock()
		heap.Pop(p)
		p.lock.Unlock()
	}
}

var tokenParser = jwt.NewParser(jwt.WithoutClaimsValidation())

func (p *ProxyClientPool) GetCachedClient(tokenString string) (ClientWithResponsesInterface, bool, error) {
	p.cleanup()
	claims := jwt.RegisteredClaims{}
	_, _, err := tokenParser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, false, err
	}
	if claims.ExpiresAt.Before(p.GetCutOff()) {
		return nil, true, nil
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, item := range p.items {
		if item.tokenID == claims.ID {
			return item.client, false, nil
		}
	}
	return nil, false, nil
}

func (p *ProxyClientPool) AddClient(client ClientWithResponsesInterface, tokenString string) error {
	p.cleanup()
	claims := jwt.RegisteredClaims{}
	_, _, err := tokenParser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return err
	}
	if claims.ExpiresAt.Before(p.GetCutOff()) {
		return errors.New("token is expired or pool is full for old tokens")
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.items) >= p.maxSize {
		if claims.ExpiresAt.Before(p.items[0].exp) || claims.ExpiresAt.Equal(p.items[0].exp) {
			p.cutoff = p.items[0].exp.Add(time.Nanosecond)
			// reject since client is too old
			return errors.New("pool is full, client token is too old")
		}
		// replace the oldest client
		p.items[0].client = client
		p.items[0].tokenID = claims.ID
		p.items[0].exp = claims.ExpiresAt.Time
		heap.Fix(p, 0)
	} else {
		heap.Push(p, &pooledClientItem{
			client:  client,
			tokenID: claims.ID,
			exp:     claims.ExpiresAt.Time,
		})
	}
	return nil
}

func (p *ProxyClientPool) GetCutOff() time.Time {
	now := time.Now()
	if p.cutoff.Before(now) {
		p.cutoff = now
	}
	return p.cutoff
}

func NewProxyClientPool(maxSize int) *ProxyClientPool {
	return &ProxyClientPool{
		items:   make([]*pooledClientItem, 0, maxSize),
		maxSize: maxSize,
	}
}
