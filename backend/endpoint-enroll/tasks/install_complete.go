package tasks

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/secret"
)

type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func InstallComplete(receiptIn io.Reader) error {
	receipt, err := readReceipt(receiptIn)
	if err != nil {
		return err
	}
	/*

		tenantID := common.MustGetenv(common.DefaultEnvVarAzureTenantId)
		clientID := uuid.MustParse(common.MustGetenv(common.DefaultEnvVarAzureClientId))
		endpointClientID := common.MustGetenv(common.DefaultEnvVarAppAzureClientId)

		templateGroupID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_TEMPLATE_GROUP_ID"))
		templateID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_TEMPLATE_ID"))
		deviceObjectID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_DEVICE_OBJECT_ID"))
		deviceLinkID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_DEVICE_LINK_ID"))
		servicePrincipalId := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_SERVICE_PRINCIPAL_ID"))
	*/
	//serviceClient, err := newServiceClientForInstall(clientID.String(), tenantID, endpointClientID)
	if err != nil {
		return err
	}

	//body := client.CertificateEnrollmentReplyFinalize{}

	// create key
	ss := secret.GetService(context.Background())
	header := JWTHeader{
		Alg: "RS256",
		Typ: "JWT",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	buf := bytes.Buffer{}
	buf.WriteString(headerEncoded)
	buf.WriteByte('.')
	buf.WriteString(receipt.JwtClaims)
	hash := sha256.Sum256(buf.Bytes())
	signature, pubkey, err := ss.RS256SignHash(hash[:], receipt.Ref.ID.String())
	if err != nil {
		return err
	}

	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	if err != nil {
		return err
	}
	fmt.Println(signatureEncoded)
	fmt.Println(pubkey.E)
	fmt.Println(base64.RawURLEncoding.EncodeToString(pubkey.N.Bytes()))

	return nil
}
