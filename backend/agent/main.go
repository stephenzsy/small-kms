package main

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
)

const DefaultEnvVarTenantID = "AZURE_TENANT_ID"
const DefaultEnvVarClientID = "AZURE_CLIENT_ID"
const DefaultEnvVarCertBundlePath = "AZURE_CERT_BUNDLE"
const DefaultEnvVarApiBaseUrl = "SMALLKMS_API_BASE_URL"
const DefaultEnvVarApiScope = "SMALLKMS_API_SCOPE"

func main() {
	// Find .env file
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("Error loading .env file: %s\n", err)
	}

	var clientID, tenantID, bundlePath, apiBaseUrl, apiScope string
	var ok bool
	if tenantID, ok = os.LookupEnv(DefaultEnvVarTenantID); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarTenantID)
	}
	if clientID, ok = os.LookupEnv(DefaultEnvVarClientID); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarClientID)
	}
	if bundlePath, ok = os.LookupEnv(DefaultEnvVarCertBundlePath); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarCertBundlePath)
	}
	if apiBaseUrl, ok = os.LookupEnv(DefaultEnvVarApiBaseUrl); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarApiBaseUrl)
	}
	if apiScope, ok = os.LookupEnv(DefaultEnvVarApiScope); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarApiScope)
	}
	bundleBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		log.Panicf("Failed to read certificate bundle: %v\n", err)
	}

	var privateKey crypto.PrivateKey
	var x509Certs []*x509.Certificate
	for block, rest := pem.Decode(bundleBytes); block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "PRIVATE KEY" {
			privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				{
					log.Panicf("Failed to parse private key: %v\n", err)
				}
			}
		} else if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Panicf("Failed to parse certificate: %v\n", err)
			}
			x509Certs = append(x509Certs, cert)
		}

	}

	creds, err := azidentity.NewClientCertificateCredential(tenantID, clientID, x509Certs, privateKey, nil)
	if err != nil {
		log.Panicf("Failed to create credential: %v\n", err)
	}
	client, err := agentclient.NewClientWithCreds(apiBaseUrl, creds, []string{apiScope}, tenantID)
	if err != nil {
		log.Panicf("Failed to create client: %v\n", err)
	}
	resp, err := client.AgentCheckIn(context.Background(), nil)
	if err != nil {
		log.Panicf("Failed to check in: %v\n", err)
	}
	log.Println("Responded with", resp.Status, "and", resp.Body, "body")
}
