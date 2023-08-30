package scep

import (
	"context"
	"encoding/asn1"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/scep/cms"
	"go.mozilla.org/pkcs7"
)

type pkiMessage struct {
	cms.SignedData
	TransactionID TransactionID
	MessageType   MessageType
}

// The following message types are defined:
type MessageType string

// Undefined message types are treated as an error.
const (
	CertRep    MessageType = "3"
	RenewalReq             = "17"
	UpdateReq              = "18"
	PKCSReq                = "19"
	CertPoll               = "20"
	GetCert                = "21"
	GetCRL                 = "22"
)

type TransactionID string

var (
	oidSCEPmessageType    = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 2}
	oidSCEPpkiStatus      = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 3}
	oidSCEPfailInfo       = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 4}
	oidSCEPsenderNonce    = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 5}
	oidSCEPrecipientNonce = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 6}
	oidSCEPtransactionID  = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 7}
)

func ParsePkiMessage(raw []byte) (*pkiMessage, error) {
	contentInfo, err := cms.ParseCMS(raw)
	if err != nil {
		return nil, err
	}

	// root is signedData
	if contentInfo.ContentTypeIDString() != cms.ObjIDSignedData {
		return nil, fmt.Errorf("invalid content type: %s", contentInfo.ContentTypeIDString())
	}
	// verify with mozilla lib
	if p7, mParseErr := pkcs7.Parse(raw); mParseErr == nil {
		if err = p7.Verify(); err != nil {
			return nil, err
		}
	}

	msgSignedData := contentInfo.SignedData()
	var tID TransactionID

	if err = msgSignedData.UnmarshalSignedAttribute(oidSCEPtransactionID, &tID); err != nil {
		return nil, err
	}

	var msgType MessageType
	if err = msgSignedData.UnmarshalSignedAttribute(oidSCEPmessageType, &msgType); err != nil {
		return nil, err
	}

	// todo verify message type

	return &pkiMessage{
		SignedData:    msgSignedData,
		TransactionID: tID,
		MessageType:   msgType,
	}, nil

}

/*
	func (s *scepServer) DecryptPKIEnvelope(msg *pkiMessage) error {
		encapsulatedContentInfo := msg.SignedData.EncapContentInfo()
		if encapsulatedContentInfo.ContentTypeIDString() != cms.ObjIDData {
			return fmt.Errorf("invalid content type in pki encapsulatedContentInfo: %s", encapsulatedContentInfo.ContentTypeIDString())
		}
		envelopdContentInfo := encapsulatedContentInfo.Data()
		if envelopdContentInfo.ContentTypeIDString() != cms.ObjIDEnvelopedData {
			return fmt.Errorf("invalid content type in pki eContent: %s", envelopdContentInfo.ContentTypeIDString())
		}
		csr, err := envelopdContentInfo.EnvelopedData().Decrypt(s.keyvaultDecrypt)
		if err != nil {
			return err
		}
		msg.pkiEnvelope, err = p7.Decrypt(cert, key)
		if err != nil {
			return err
		}

}
*/
func (s *scepServer) keyvaultDecrypt(ctx context.Context, kid azkeys.ID, cipherText []byte) ([]byte, error) {
	resp, err := s.azKeysClient.UnwrapKey(ctx, kid.Name(), kid.Version(), azkeys.KeyOperationParameters{
		Value:     cipherText,
		Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSA15),
	}, nil)
	return resp.Result, err
	/*
		data, ok := p7.raw.(envelopedData)
		if !ok {
			return nil, ErrNotEncryptedContent
		}
		recipient := selectRecipientForCertificate(data.RecipientInfos, cert)
		if recipient.EncryptedKey == nil {
			return nil, errors.New("pkcs7: no enveloped recipient for provided certificate")
		}
		switch pkey := pkey.(type) {
		case *rsa.PrivateKey:
			var contentKey []byte
			contentKey, err := rsa.DecryptPKCS1v15(rand.Reader, pkey, recipient.EncryptedKey)
			if err != nil {
				return nil, err
			}
			return data.EncryptedContentInfo.decrypt(contentKey)
		}
		return nil, ErrUnsupportedAlgorithm
	*/
}
