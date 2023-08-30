package scep

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/admin"
)

const defaultCaps = "Renewal\nSHA-256\nAES\nSCEPStandard\nPOSTPKIOperation"

func (s *scepServer) refreshCaCert(c *gin.Context) error {
	resp, err := s.smallKmsAdminClient.ListCertificatesV1(c, uuid.Nil) // TODO
	if err != nil {
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Faild to get certificate metadata: %s", err.Error())
		c.Data(500, "text/plain", []byte("Internal Error"))
		return err
	}
	certRefs := []admin.CertificateRef{}
	if err = json.Unmarshal(respBody, &certRefs); err != nil {
		return err
	}
	if len(certRefs) == 0 {
		return errors.New("no certificated available")
	}

	getCertResp, err := s.smallKmsAdminClient.GetCertificateV1(c, certRefs[0].NamespaceID, certRefs[0].ID, &admin.GetCertificateV1Params{
		Accept: to.Ptr(admin.AcceptX509CaCert),
	})
	if err != nil {
		return err
	}
	getCertResp.Header.Get("X-Keyvault-Kid")
	blob, err := io.ReadAll(getCertResp.Body)
	if err != nil {
		return err
	}
	s.cachedCaCert = blob
	return nil
}

func (s *scepServer) HandleGetCaCert(c *gin.Context) {
	if len(s.cachedCaCert) == 0 || s.cachedCaTime.Add(time.Hour*24).Before(time.Now()) {
		s.cachedCaTime = time.Now()
		if err := s.refreshCaCert(c); err != nil {
			s.cachedCaCert = nil
			log.Printf("Faild to get certificate: %s", err.Error())
			c.Data(500, "text/plain", []byte("Internal Error"))
			return
		}
	}
	c.Data(200, "application/x-x509-ca-cert", s.cachedCaCert)
}

const maxPayloadSize = 2 << 20

func (s *scepServer) ScepGet(c *gin.Context, params ScepGetParams) {
	switch params.Operation {
	case GetCACaps:
		c.Data(200, "text/plain", []byte(defaultCaps))
		return
	case GetCaCert:
		s.HandleGetCaCert(c)
		return
	}
	c.Data(404, "text/plain", []byte("Not Found"))
}

func (s *scepServer) ScepPost(c *gin.Context, params ScepPostParams) {
	switch params.Operation {
	case PKIOperation:
		raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxPayloadSize))
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
