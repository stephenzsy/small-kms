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

func (identifier *CertificateIdentifier) docID() kmsdoc.KmsDocID {
	if identifier == nil {
		return kmsdoc.NewKmsDocID(kmsdoc.DocTypeUnknown, uuid.Nil)
	}
	if identifier.Type != nil {
		switch *identifier.Type {
		case CertIdTypePolicyId:
			return kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, identifier.ID)
		}
	}
	return kmsdoc.NewKmsDocID(kmsdoc.DocTypeCert, identifier.ID)
}

func docIDtoCertIdentifier(id kmsdoc.KmsDocID) CertificateIdentifier {
	switch id.GetType() {
	case kmsdoc.DocTypeLatestCertForPolicy:
		return CertificateIdentifier{
			ID:   id.GetUUID(),
			Type: ToPtr(CertIdTypePolicyId),
		}
	case kmsdoc.DocTypeCert:
		return CertificateIdentifier{
			ID:   id.GetUUID(),
			Type: ToPtr(CertIdTypeCertId),
		}
	}
	return CertificateIdentifier{
		ID: id.GetUUID(),
	}
}

func (s *adminServer) getCertDoc(c context.Context, namespaceID uuid.UUID, id kmsdoc.KmsDocID) (*CertDoc, error) {
	pd := new(CertDoc)
	err := kmsdoc.AzCosmosRead(c, s.azCosmosContainerClientCerts, namespaceID,
		id, pd)
	return pd, err
}

func (d *CertDoc) GetCUID() kmsdoc.KmsDocID {
	if d.AliasID != nil {
		return *d.AliasID
	}
	return d.ID
}

/*

func (s *adminServer) listCertDocForPolicyID(ctx context.Context, namespaceID uuid.UUID, policyID uuid.UUID) ([]*CertDoc, error) {

	partitionKey := azcosmos.NewPartitionKeyString(namespaceID.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.odType,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c.type = @type`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: directoryID.String()},
				{Name: "@type", Value: kmsdoc.DocTypeNameCert},
			},
		})

	return PagerToList[DirectoryObjectDoc](ctx, pager)
}
*/
