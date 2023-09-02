package cms

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
)

type CMSContentTypeIDStr string

// Object identifier strings of the three implemented PKCS7 types.
const (
	ObjIDData          CMSContentTypeIDStr = "1.2.840.113549.1.7.1"
	ObjIDSignedData    CMSContentTypeIDStr = "1.2.840.113549.1.7.2"
	ObjIDEnvelopedData CMSContentTypeIDStr = "1.2.840.113549.1.7.3"
)

var (
	oidData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	oidSignedData    = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	oidEnvelopedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 3}

	oidSCEPtransactionID = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 7}
	oidSCEPmessageType   = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 2}
	oidSCEPsenderNonce   = asn1.ObjectIdentifier{2, 16, 840, 1, 113733, 1, 9, 5}

	// Encryption Algorithms
	oidEncryptionAlgorithmAES128CBC = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 2}
)

const (
	MessageTypePKCSReq    = "19"
	MessageTypeRenewalReq = "17"
)

// RFC 5652 Section 3
type contentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"tag:0,explicit,optional"`
}
type contentInfoParsed struct {
	contentInfo
	dataParsed          *contentInfoParsed
	signedDataParsed    *signedDataParsed
	envelopedDataParsed *envelopedDataParsed
}
type ContentInfo interface {
	ContentTypeIDString() CMSContentTypeIDStr
	Data() ContentInfo
	SignedData() SignedData
	EnvelopedData() EnvelopedData
}

func (ci *contentInfo) ContentTypeIDString() CMSContentTypeIDStr {
	return CMSContentTypeIDStr(ci.ContentType.String())
}
func (ci *contentInfoParsed) Data() ContentInfo {
	return ci.dataParsed
}
func (ci *contentInfoParsed) SignedData() SignedData {
	return ci.signedDataParsed
}
func (ci *contentInfoParsed) EnvelopedData() EnvelopedData {
	return ci.envelopedDataParsed
}

// RFC 5652 Section 5.1
type signedData struct {
	Version          int
	DigestAlgorithms asn1.RawValue
	EncapContentInfo contentInfo
	Certificates     asn1.RawValue `asn1:"optional,tag:0"`
	Crls             asn1.RawValue `asn1:"optional,tag:1"`
	SignerInfos      []signerInfo  `asn1:"set"`
}
type signedDataParsed struct {
	signedData
	certificates     []*x509.Certificate
	crl              *x509.RevocationList
	encapContentInfo *contentInfoParsed
}
type SignedData interface {
	UnmarshalSignedAttribute(attributeType asn1.ObjectIdentifier, out any) error
	EncapContentInfo() ContentInfo
}

func (s *signedData) toParsed() (t signedDataParsed, err error) {
	t.signedData = *s
	if len(s.Certificates.Bytes) > 0 {
		if t.certificates, err = x509.ParseCertificates((s.Certificates.Bytes)); err != nil {
			return
		}
	}
	if len(s.Crls.Bytes) != 0 {
		if t.crl, err = x509.ParseRevocationList(s.Crls.Bytes); err != nil {
			return
		}
	}
	if len(s.SignerInfos) == 0 {
		err = errors.New("pkcs7: no signer infos")
		return
	}

	var compound asn1.RawValue
	var content []byte

	// The Content.Bytes maybe empty on PKI responses.
	if len(s.EncapContentInfo.Content.Bytes) > 0 {
		if _, err = asn1.Unmarshal(s.EncapContentInfo.Content.Bytes, &compound); err != nil {
			return
		}
	}
	// Compound octet string
	if compound.IsCompound {
		if compound.Tag == 4 {
			if _, err = asn1.Unmarshal(compound.Bytes, &content); err != nil {
				return
			}
		} else {
			content = compound.Bytes
		}
	} else {
		// assuming this is tag 04
		content = compound.Bytes
	}

	if t.encapContentInfo, err = ParseCMS(content); err != nil {
		err = fmt.Errorf("parse inner content:\n%w", err)
		return
	}
	return
}

func (sd *signedData) UnmarshalSignedAttribute(attributeType asn1.ObjectIdentifier, out any) error {
	return sd.SignerInfos[0].UnmarshalSignedAttribute(attributeType, out)
}
func (sd *signedDataParsed) EncapContentInfo() ContentInfo {
	return sd.encapContentInfo
}

// RFC 5652 Section 5.3
type signerInfo struct {
	Version               int
	IssuerAndSerialNumber issuerAndSerialNumber
	DigestAlgorithm       asn1.RawValue
	SignedAttrs           []attribute `asn1:"set,optional,omitempty,tag:0"`
	SignatureAlgorithm    asn1.RawValue
	Signature             asn1.RawValue
	UnsignedAttrs         asn1.RawValue `asn1:"optional,omitempty,tag:1"`
}
type attribute struct {
	AttrType   asn1.ObjectIdentifier
	AttrValues asn1.RawValue `asn1:"set"`
}

// RFC 5652 Section 6.1
type envelopedData struct {
	Version              int
	OriginatorInfo       asn1.RawValue   `asn1:"optional,omitempty,tag:0"`
	RecipientInfos       []recipientInfo `asn1:"set"`
	EncryptedContentInfo encryptedContentInfo
	UnprotectedAttrs     asn1.RawValue `asn1:"optional,omitempty,tag:1"`
}
type encryptedContentInfo struct {
	ContentType                asn1.ObjectIdentifier
	ContentEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedContent           asn1.RawValue `asn1:"tag:0,optional"`
}
type envelopedDataParsed struct {
	envelopedData
}
type EnvelopedData interface {
	Decrypt(cert *x509.Certificate, keyUnwrapper KeyUnwrapperRSA1_5) ([]byte, error)
}

// RFC 5652 Section 6.2, 6.2.1 KeyTransRecipientInfo
type recipientInfo struct {
	Version                int
	IssuerAndSerialNumber  issuerAndSerialNumber
	KeyEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedKey           []byte
}

// RFC 5652 Section 10.2.4
type issuerAndSerialNumber struct {
	IssuerName   asn1.RawValue
	SerialNumber *big.Int
}

func (si *signerInfo) UnmarshalSignedAttribute(attributeType asn1.ObjectIdentifier, out any) error {
	for _, attr := range si.SignedAttrs {
		if attr.AttrType.Equal(attributeType) {
			_, err := asn1.Unmarshal(attr.AttrValues.Bytes, out)
			return err
		}
	}
	return fmt.Errorf("cms: attribute type %s not in attributes", attributeType.String())
}

func (s *envelopedData) toParsed() (t envelopedDataParsed, err error) {
	t.envelopedData = *s
	alg := s.EncryptedContentInfo.ContentEncryptionAlgorithm.Algorithm
	if !alg.Equal(oidEncryptionAlgorithmAES128CBC) {
		err = fmt.Errorf("unsupported content encryption algorithm: %s", alg)
	}
	return
}

func isCertMatchForIssuerAndSerial(cert *x509.Certificate, ias issuerAndSerialNumber) bool {
	return cert.SerialNumber.Cmp(ias.SerialNumber) == 0 &&
		bytes.Equal(cert.RawIssuer, ias.IssuerName.FullBytes)
}

func selectRecipientForCertificate(recipients []recipientInfo, cert *x509.Certificate) recipientInfo {
	for _, recp := range recipients {
		if isCertMatchForIssuerAndSerial(cert, recp.IssuerAndSerialNumber) {
			return recp
		}
	}
	return recipientInfo{}
}

type KeyUnwrapperRSA1_5 func(content []byte) ([]byte, error)

func (s *envelopedDataParsed) Decrypt(cert *x509.Certificate, keyUnwrapper KeyUnwrapperRSA1_5) ([]byte, error) {
	recipient := selectRecipientForCertificate(s.RecipientInfos, cert)

	if recipient.EncryptedKey == nil {
		return nil, errors.New("pkcs7: no enveloped recipient for provided certificate")
	}

	var contentKey []byte
	contentKey, err := keyUnwrapper(recipient.EncryptedKey)
	if err != nil {
		return nil, err
	}
	return s.EncryptedContentInfo.decrypt(contentKey)
}

func (eci *encryptedContentInfo) decrypt(key []byte) ([]byte, error) {
	alg := eci.ContentEncryptionAlgorithm.Algorithm

	// EncryptedContent can either be constructed of multple OCTET STRINGs
	// or _be_ a tagged OCTET STRING
	var cyphertext []byte
	if eci.EncryptedContent.IsCompound {
		// Complex case to concat all of the children OCTET STRINGs
		var buf bytes.Buffer
		cypherbytes := eci.EncryptedContent.Bytes
		for {
			var part []byte
			cypherbytes, _ = asn1.Unmarshal(cypherbytes, &part)
			buf.Write(part)
			if cypherbytes == nil {
				break
			}
		}
		cyphertext = buf.Bytes()
	} else {
		// Simple case, the bytes _are_ the cyphertext
		cyphertext = eci.EncryptedContent.Bytes
	}

	var block cipher.Block
	var err error

	switch {
	case alg.Equal(oidEncryptionAlgorithmAES128CBC):
		block, err = aes.NewCipher(key)
	}

	if err != nil {
		return nil, err
	}

	iv := eci.ContentEncryptionAlgorithm.Parameters.Bytes
	if len(iv) != block.BlockSize() {
		return nil, errors.New("pkcs7: encryption algorithm parameters are malformed")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(cyphertext))
	mode.CryptBlocks(plaintext, cyphertext)
	if plaintext, err = unpad(plaintext, mode.BlockSize()); err != nil {
		return nil, err
	}
	return plaintext, nil
}

func unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen < 1 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}

	// the last byte is the length of padding
	padlen := int(data[len(data)-1])

	// check padding integrity, all bytes should be the same
	pad := data[len(data)-padlen:]
	for _, padbyte := range pad {
		if padbyte != byte(padlen) {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-padlen], nil
}

// parse only types related to pki message
func ParseCMS(raw []byte) (r *contentInfoParsed, err error) {
	/*
		r.p7, err = pkcs7.Parse(raw)
		if err != nil {
			return
		}

		// expect p7 is Signed data
		if err = r.p7.Verify(); err != nil {
			return
		}

		if err = r.p7.UnmarshalSignedAttribute(oidSCEPtransactionID, &r.transactionId); err != nil {
			return
		}

		if err = r.p7.UnmarshalSignedAttribute(oidSCEPmessageType, &r.messageType); err != nil {
			return
		}

		switch r.messageType {
		case MessageTypePKCSReq, MessageTypeRenewalReq:

			if err = r.p7.UnmarshalSignedAttribute(oidSCEPsenderNonce, &r.senderNonce); err != nil {
				return
			}
			if err = parseEnvelopedData(r.p7.Content, &r.envelope); err != nil {
				return
			}
		default:
			err = fmt.Errorf("unsupported message type: %s", r.messageType)
			return
		}
		return
	*/

	var p7 contentInfo
	der, err := ber2der(raw)
	if err != nil {
		return nil, fmt.Errorf("ber2der: %s", err)
	}

	if _, err := asn1.Unmarshal(der, &p7); err != nil {
		return nil, fmt.Errorf("ansi unmarshall error:\n%s", err)
	}
	cip := contentInfoParsed{}
	switch {
	case oidData.Equal(p7.ContentType):
		dataParsed, err := ParseCMS(p7.Content.Bytes)
		if err != nil {
			return nil, err
		}
		cip.dataParsed = dataParsed
	case oidSignedData.Equal(p7.ContentType):
		signedDataRaw := signedData{}
		if _, err := asn1.Unmarshal(p7.Content.Bytes, &signedDataRaw); err != nil {
			return nil, fmt.Errorf("unmarshal signed data:\n%w", err)
		}
		signedData, err := signedDataRaw.toParsed()
		if err != nil {
			return nil, fmt.Errorf("parse signed data:\n%w", err)
		}
		cip.signedDataParsed = &signedData
	case oidEnvelopedData.Equal(p7.ContentType):
		envelopedDataRaw := envelopedData{}
		if _, err := asn1.Unmarshal(p7.Content.Bytes, &envelopedDataRaw); err != nil {
			return nil, err
		}
		envelopedData, err := envelopedDataRaw.toParsed()
		if err != nil {
			return nil, err
		}
		cip.envelopedDataParsed = &envelopedData
	default:
		return nil, fmt.Errorf("unsupported content type: %s", p7.ContentType.String())
	}
	return &cip, nil
}

type ReqPkiMessage struct {
	cms           ContentInfo
	TransactionId string
	MessageType   string
	SenderNonce   []byte
}

// parse only types related to pki message
func ParsePkiMessage(raw []byte) (r ReqPkiMessage, err error) {
	r.cms, err = ParseCMS(raw)
	if err != nil {
		return
	}
	sd := r.cms.SignedData()
	if err = sd.UnmarshalSignedAttribute(oidSCEPtransactionID, &r.TransactionId); err != nil {
		return
	}
	if err = sd.UnmarshalSignedAttribute(oidSCEPmessageType, &r.MessageType); err != nil {
		return
	}
	if err = sd.UnmarshalSignedAttribute(oidSCEPsenderNonce, &r.SenderNonce); err != nil {
		return
	}
	return
}

func (r *ReqPkiMessage) Decrypt(cert *x509.Certificate, keyUnwrapper KeyUnwrapperRSA1_5) ([]byte, error) {
	return r.cms.SignedData().EncapContentInfo().EnvelopedData().Decrypt(cert, keyUnwrapper)
}
