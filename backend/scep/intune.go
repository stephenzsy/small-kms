package scep

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/scep/cms"
	"github.com/stephenzsy/small-kms/backend/scep/msintune"
)

/*

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
*/

const DEFAULT_MSGRAPH_RESOURCE_URL = "https://graph.microsoft.com/"
const INTUNE_API_SCOPE = "https://api.manage.microsoft.com/.default"
const DEFAULT_INTUNE_APP_ID = "0000000a-0000-0000-c000-000000000000"
const VALIDATION_SERVICE_NAME = "ScepRequestValidationFEService"

func newIntuneClient(creds azcore.TokenCredential, endpoint string) (client *msintune.Client, err error) {
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{INTUNE_API_SCOPE},
	}
	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Renegotiation: tls.RenegotiateFreelyAsClient,
				MinVersion:    tls.VersionTLS12,
				// You may need this if connecting to servers with self-signed certificates
				// InsecureSkipVerify: true,
			},
		},
	}
	client, err = msintune.NewClient(endpoint, msintune.WithHTTPClient(&httpClient), msintune.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Client-Request-Id", uuid.New().String())
		req.Header.Add("Api-Version", "2018-02-20")
		tokenResp, err := creds.GetToken(ctx, tokenOptions)
		if err != nil {
			return fmt.Errorf("failed to get creds: %w", err)
		}
		req.Header.Add("Authorization", "Bearer "+tokenResp.Token)
		fmt.Println(req.Header)
		return nil
	}))
	return
}

func (s *scepServer) refreshServiceMap() (client *msintune.Client, err error) {
	sp, err := s.msGraphClient.ServicePrincipalsWithAppId(to.Ptr(DEFAULT_INTUNE_APP_ID)).Get(context.Background(), nil)
	if err != nil {
		return
	}
	endpoints := sp.GetEndpoints()
	for _, endpoint := range endpoints {
		if *endpoint.GetProviderName() == VALIDATION_SERVICE_NAME {
			s.msIntuneScepEndpoint = *endpoint.GetUri()
		}
	}
	return newIntuneClient(s.DefaultAzCredential(), s.msIntuneScepEndpoint)
}

func (s *scepServer) getIntuneClient() (client *msintune.Client, err error) {
	if s.msIntuneClient == nil || s.rateLimiters.msGraphServiceMapping.Allow() {
		s.msIntuneClient, err = s.refreshServiceMap()
		if err != nil {
			s.msIntuneClient = nil
		}
	}
	return s.msIntuneClient, err
}

func validatePkiMessageWithIntune(ctx context.Context, intuneClient *msintune.Client, csrDer []byte, pkiMessage cms.ReqPkiMessage) error {
	reqBody := msintune.ValidateRequestJSONRequestBody{
		Request: msintune.IntuneScepValidateRequest{
			TransactionId:      pkiMessage.TransactionId,
			CertificateRequest: csrDer,
			CallerInfo:         to.Ptr("test"),
		},
	}
	resp, err := intuneClient.ValidateRequest(ctx, reqBody)
	if err != nil {
		return err
	}
	reqJson, err := json.Marshal(reqBody)
	fmt.Println(string(reqJson))

	bytes, err := io.ReadAll(resp.Body)
	return fmt.Errorf("body: %s\n%v", string(bytes), resp)
}
