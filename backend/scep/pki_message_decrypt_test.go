package scep

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/joho/godotenv"
	"github.com/stephenzsy/small-kms/backend/scep/cms"
)

func TestCertDecrypt(t *testing.T) {
	//t.Skip()
	c, err := os.ReadFile("cms/testdata/cacert.pem")
	if err != nil {
		t.Skip(err)
	}
	b, err := os.ReadFile("cms/testdata/pki-payload.txt")
	if err != nil {
		t.Skip(err)
	}
	payload, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		t.Error(err)
	}
	cp, _ := pem.Decode(c)
	if cp == nil {
		t.Skip("No cert tro select")
	}
	cacert, err := x509.ParseCertificate(cp.Bytes)
	if err != nil {
		t.Error(err)
	}

	msg, err := cms.ParsePkiMessage(payload)
	if err != nil {
		t.Error(err)
	}
	err = godotenv.Load("cms/testdata/.env.scep")
	if err != nil {
		t.Skip(err)
	}
	azCreds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Error(err)
	}
	azKeysClient, err := azkeys.NewClient(os.Getenv("AZURE_KEYVAULT_RESOURCEENDPOINT"), azCreds, nil)
	if err != nil {
		t.Error(err)
	}
	csr, err := msg.Decrypt(cacert, func(wrappedKey []byte) ([]byte, error) {
		resp, err := azKeysClient.UnwrapKey(context.Background(), os.Getenv("TEST_KEYVAULT_KEY_NAME"), os.Getenv("TEST_KEYVAULT_KEY_VERSION"), azkeys.KeyOperationParameters{
			Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSA15),
			Value:     wrappedKey,
		}, nil)
		return resp.Result, err
	})
	os.WriteFile("cms/testdata/csr.der", csr, 0644)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = x509.ParseCertificateRequest(csr)
	if err != nil {
		t.Error(err)
	}

	client, err := newIntuneClient(azCreds, os.Getenv("INTUNE_SCEP_ENDPOINT"))
	if err != nil || client == nil {
		t.Error(err)
		return
	}

	if err = validatePkiMessageWithIntune(context.Background(), client, csr, msg); err != nil {
		t.Error(err.Error())
		t.Error(err)
	}

}
