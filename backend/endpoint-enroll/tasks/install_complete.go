package tasks

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/client"
	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/secret"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func InstallComplete(receiptIn io.Reader, installToUser bool) error {
	receipt, err := readReceipt(receiptIn)
	if err != nil {
		return err
	}

	tenantID := common.MustGetenv(common.IdentityEnvVarNameAzTenantID)
	clientID := uuid.MustParse(common.MustGetenv(common.IdentityEnvVarNameAzClientID))
	endpointClientID := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientID, "")

	serviceClient, err := newServiceClientForInstall(clientID.String(), tenantID, endpointClientID)
	if err != nil {
		return err
	}

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
	installToMachine := !installToUser
	signature, pubkey, err := ss.RS256SignHash(hash[:], receipt.Ref.ID.String(), installToMachine)
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

	body := client.CertificateEnrollmentReplyFinalize{
		JwtHeader:    headerEncoded,
		JwtSignature: signatureEncoded,
	}
	//rsaPublicKeyPopulateJwk(pubkey, &body.PublicKey)
	serviceClient.CompleteCertificateEnrollmentV2(context.Background(), receipt.Ref.NamespaceID, receipt.Ref.ID, nil, body)

	return nil
}

func rsaPublicKeyPopulateJwk(pubkey *rsa.PublicKey, p *shared.JwkProperties) {
	if pubkey == nil {
		return
	}

	p.Alg = utils.ToPtr(shared.AlgRS256)
	p.Kty = "RSA"
	p.E = big.NewInt(int64(pubkey.E)).Bytes()
	p.N = pubkey.N.Bytes()
	p.KeySize = utils.ToPtr(int32(pubkey.Size()) * 8)
}
