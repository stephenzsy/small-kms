package scep

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/scep/cms"
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

func (s *scepServer) handlePkiMessageRequest(ctx context.Context, pkiMessage cms.ReqPkiMessage, namespaceID uuid.UUID) error {
	item := s.getCaCert(ctx, namespaceID)
	if item.status != http.StatusOK {
		return errors.New("failed to get ca cert")
	}

	kid := azkeys.ID(item.certItem.KeyStore)

	csrDer, err := pkiMessage.Decrypt(item.cert, func(wrappedKey []byte) ([]byte, error) {
		resp, err := s.AzKeysClient().UnwrapKey(context.Background(), kid.Name(), kid.Version(), azkeys.KeyOperationParameters{
			Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSA15),
			Value:     wrappedKey,
		}, nil)
		return resp.Result, err
	})
	if csrDer == nil || err != nil {
		return fmt.Errorf("failed to decrypt request: %w", err)
	}

	intuneClient, err := s.getIntuneClient()
	if err != nil {
		return err
	}
	if intuneClient == nil {
		return errors.New("intune client is nil")
	}

	validatePkiMessageWithIntune(ctx, intuneClient, csrDer, pkiMessage)
	/*

		csrDerBase64 := base64.StdEncoding.EncodeToString(csrDer)
		validateReq := models.NewIntuneScepValidateRequest()
		validateReq.SetTransactionId(&pkiMessage.TransactionId)
		validateReq.SetCertificateRequest(&csrDerBase64)
		validateReq.SetCallerInfo(to.Ptr("SmallKMS:0.0.1"))
		validateBody := scepactions.NewValidateRequestPostRequestBody()
		validateBody.SetRequest(validateReq)
		resp, err := intuneClient.ScepActions().ValidateRequest().Post(ctx, validateBody, nil)
		if err != nil {
			return err
		}

		fmt.Printf("intunes response: %v\n", resp)
	*/
	// send to intune
	return nil
}

func (s *scepServer) ScepPost(c *gin.Context, namespaceID uuid.UUID, params ScepPostParams) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxPayloadSize))
	log.Printf("Log: %s\n%s\n", params.Operation, base64.StdEncoding.EncodeToString(raw))
	if err != nil {
		c.Data(400, "text/plain", nil)
		return
	}
	switch params.Operation {
	case PKIOperation:

		pkiMessage, err := cms.ParsePkiMessage(raw)
		if err != nil {
			c.Data(400, "text/plain", []byte("Unable to parse"))
			return
		}
		switch pkiMessage.MessageType {
		case cms.MessageTypePKCSReq,
			cms.MessageTypeRenewalReq:
			err = s.handlePkiMessageRequest(c, pkiMessage, namespaceID)
			if err != nil {
				c.Data(500, "text/plain", nil)
				return
			}
			return
		}
	}

	c.Data(400, "text/plain", nil)
}
