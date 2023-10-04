package tasks

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/client"
)

func readReceipt(reader io.Reader) (*client.CertificateEnrollmentReceipt, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	receipt := client.CertificateEnrollmentReceipt{}
	err = json.Unmarshal(data, &receipt)
	return &receipt, err
}

func InstallView(receiptIn io.Reader) error {
	receipt, err := readReceipt(receiptIn)
	if err != nil {
		return err
	}
	fmt.Printf("Certificate enrollment receipt\n")
	fmt.Printf("==============================\n")
	fmt.Printf("Must be completed before \033[0;33m%s\033[0m\n", receipt.Expires)
	keyInfoJson, _ := json.MarshalIndent(receipt.KeyProperties, "", "  ")
	fmt.Printf("Key info: %s\n", keyInfoJson)
	jwtClaimsDecoded, err := base64.RawURLEncoding.DecodeString(receipt.JwtClaims)
	if err != nil {
		return err
	}
	claimsJson := bytes.Buffer{}
	json.Indent(&claimsJson, jwtClaimsDecoded, "", "  ")
	fmt.Printf("JWT claims: %s\n", claimsJson.String())
	return nil
}
