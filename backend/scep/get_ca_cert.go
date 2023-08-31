package scep

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/admin"
)

type caCertCachedItem struct {
	certItem admin.CertDBItem
	cert     *x509.Certificate
	status   int
}

type getCaCertCache struct {
	mu   sync.Mutex
	last time.Time

	item caCertCachedItem
}

func (s *scepServer) refreshCaCert(ctx context.Context, namespaceID uuid.UUID) (item caCertCachedItem, err error) {
	item.status = http.StatusInternalServerError
	policy, err := s.adminServer.ReadCertEnrollPolicyDBItem(ctx, namespaceID)
	if err != nil {
		return
	}
	if len(policy.PolicyID) == 0 {
		// no policy found
		item.status = http.StatusNotFound
		return
	}

	item.certItem, err = s.adminServer.ReadCertDBItem(ctx, namespaceID, policy.IssuerID)
	if err != nil {
		return
	}
	if item.certItem.ID == uuid.Nil {
		item.status = http.StatusNotFound
		return
	}

	pemBlob, err := s.adminServer.FetchCertificatePEMBlob(ctx, item.certItem.CertStore)
	if err != nil {
		return
	}
	p, _ := pem.Decode(pemBlob)
	if p == nil {
		return item, errors.New("failed to decode pem")
	}
	item.cert, err = x509.ParseCertificate(p.Bytes)
	if err != nil {
		return item, errors.New("failed to parse certificate")
	}
	item.status = http.StatusOK
	return
}

func (s *scepServer) getCaCert(ctx context.Context, namespaceID uuid.UUID) caCertCachedItem {
	nsCache, ok := s.getCaCertCaches[namespaceID]
	if !ok {
		nsCache = new(getCaCertCache)
		nsCache.item.status = http.StatusInternalServerError
		s.getCaCertCaches[namespaceID] = nsCache
	}
	nsCache.mu.Lock()
	defer nsCache.mu.Unlock()
	if time.Since(nsCache.last) > time.Hour || nsCache.item.status >= 500 {
		refreshed, err := s.refreshCaCert(ctx, namespaceID)
		if err != nil {
			log.Printf("Faild to get certificate: %s", err.Error())
		}
		if refreshed.status < 500 {
			nsCache.item = refreshed
			nsCache.last = time.Now()
		}
	}
	return nsCache.item
}

func (s *scepServer) HandleGetCaCert(c *gin.Context, namespaceID uuid.UUID) {
	// validate namespace id
	if namespaceID != intranetNamespaceID {
		c.Data(http.StatusNotFound, "text/plain", nil)
		return
	}

	item := s.getCaCert(c, namespaceID)
	if item.status >= 400 {
		c.Data(item.status, "text/plain", nil)
		return
	}
	c.Data(item.status, "application/x-x509-ca-cert", item.cert.Raw)
}
