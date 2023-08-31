package scep

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const defaultCaps = "Renewal\nSHA-256\nAES\nSCEPStandard\nPOSTPKIOperation"

const maxPayloadSize = 2 << 20

func (s *scepServer) ScepGet(c *gin.Context, namespaceID uuid.UUID, params ScepGetParams) {
	switch params.Operation {
	case GetCACaps:
		c.Data(200, "text/plain", []byte(defaultCaps))
		return
	case GetCACert:
		if !s.rateLimiters.getCa.Allow() {
			c.Data(http.StatusTooManyRequests, "text/plain", []byte("Too Many Requests"))
			return
		}
		s.HandleGetCaCert(c, namespaceID)
		return
	}
	c.Data(404, "text/plain", nil)
}

func (s *scepServer) ScepPost(c *gin.Context, namespaceID uuid.UUID, params ScepPostParams) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxPayloadSize))
	log.Printf("Log: %s\n%s\n", params.Operation, base64.StdEncoding.EncodeToString(raw))
	c.Data(404, "text/plain", []byte("Not Found"))
	switch params.Operation {
	case PKIOperation:
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal error"))
			return
		}
		_, err = ParsePkiMessage(raw)
		if err != nil {
			c.Data(400, "text/plain", []byte("Unable to parse"))
			return
		}
		/*
			if err := s.DecryptPKIEnvelope(msg); err != nil {
				c.Data(500, "text/plain", []byte("Internal error"))
				return
			}
		*/
	}

	c.Data(404, "text/plain", []byte("Not Found"))
}
