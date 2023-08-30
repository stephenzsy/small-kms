package scep

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	khttp "github.com/microsoft/kiota-http-go"
)

type ValidateRequest struct {
	TransactionId      uuid.UUID `json:"transactionId"`
	CertificateRequest []byte    `json:"certificateRequest"`
	CallerInfo         string    `json:"callerInfo"`
}

type ValidateRequestBody struct {
	Request ValidateRequest `json:"request"`
}

type SuccessNotification struct {
	ValidateRequest
	CertificateThumbprint         string `json:"certificateThumbprint"`
	CertificateSerialNumber       string `json:"certificateSerialNumber"`
	CcertificateExpirationDateUtc string `json:"certificateExpirationDateUtc"`
	IssuingAuthority              string `json:"issuingCertificateAuthority"`
	CaConfiguration               string `json:"caConfiguration"`
	CertificateAuthority          string `json:"certificateAuthority"`
}

type FailureNotification struct {
	ValidateRequest
	HResult          int64  `json:"hResult"`
	ErrorDescription string `json:"errorDescription"`
}

const DEFAULT_MSGRAPH_RESOURCE_URL = "https://graph.microsoft.com/"
const DEFAULT_INTUNE_APP_ID = "0000000a-0000-0000-c000-000000000000"
const VALIDATION_SERVICE_NAME = "ScepRequestValidationFEService"

func (s *scepServer) refreshServiceMap() error {
	sp, err := s.msGraphClient.ServicePrincipalsWithAppId(to.Ptr(DEFAULT_INTUNE_APP_ID)).Get(context.Background(), nil)
	if err != nil {
		return err
	}
	endpoints := sp.GetEndpoints()
	for _, endpoint := range endpoints {
		if *endpoint.GetProviderName() == VALIDATION_SERVICE_NAME {
			s.msIntunesScepEndpoint = *endpoint.GetUri()
			s.msIntunesScepEndpointRefresh = time.Now()
		}
	}
	_ = khttp.GetDefaultClient()

	return nil
}
