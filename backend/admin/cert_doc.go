package admin

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertDoc struct {
	kmsdoc.BaseDoc

	// associated policy id
	PolicyID uuid.UUID `json:"policyId"`
	Expires  time.Time `json:"expires"`
	// alias for certs with L prefix
	AliasID         *kmsdoc.KmsDocID `json:"aliasId,omitempty"`
	Usage           CertificateUsage `json:"usage"`
	IssuerNamespace uuid.UUID        `json:"issuerNamespace"`
	IssuerID        kmsdoc.KmsDocID  `json:"issuerId"`
	// display name of the certificate
	Name string `json:"name"`
	// keyvault certificate id
	CID string `json:"cid"`
	// keyvault certificate key id
	KID string `json:"kid"`
	// keyvault certificate secret id
	SID string `json:"sid"`
	// certificate storage path in blob storage
	CertStorePath string `json:"certStorePath"`
}

func (s *adminServer) GetLatestCertDocForPolicy(c context.Context, namespaceID uuid.UUID, policyID uuid.UUID) (*CertDoc, error) {
	pd := new(CertDoc)
	err := kmsdoc.AzCosmosRead(c, s.azCosmosContainerClientCerts, namespaceID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, policyID), pd)
	return pd, err
}
