package admin

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyDocSectionIssuerProperties struct {
	IssuerNamespaceID uuid.UUID `json:"issuerNamespaceId"`
	IssuerPolicyID    uuid.UUID `json:"issuerPolicyId"`
}

func (t *PolicyDocSectionIssuerProperties) validateAndFillWithCertRequestParameters(p *CertificateRequestPolicyParameters) (err error) {
	t.IssuerNamespaceID = p.IssuerNamespaceID

	if p.IssuerPolicyIdentifier == nil {
		t.IssuerPolicyID = defaultPolicyIdCertRequest
	} else if t.IssuerPolicyID, err = resolvePolicyIdentifier(*p.IssuerPolicyIdentifier); err != nil {
		return err
	}
	return nil
}

type PolicyDoc struct {
	kmsdoc.BaseDoc
	Enabled        bool                            `json:"enabled"`
	PolicyType     PolicyType                      `json:"policyType"`
	CertRequest    *PolicyCertRequestDocSection    `json:"certRequest,omitempty"`
	CertEnroll     *PolicyCertEnrollDocSection     `json:"certEnroll,omitempty"`
	CertAadAppCred *PolicyCertAadAppCredDocSection `json:"certAadAppCred,omitempty"`
}

func (s *adminServer) GetPolicyDoc(c context.Context, namespaceID uuid.UUID, policyID uuid.UUID) (*PolicyDoc, error) {
	pd := new(PolicyDoc)
	err := kmsdoc.AzCosmosRead(c, s.azCosmosContainerClientCerts, namespaceID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicy, policyID), pd)
	return pd, err
}

func (s *adminServer) deletePolicyDoc(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID, purge bool) error {
	return kmsdoc.AzCosmosDelete(c, s.azCosmosContainerClientCerts, namespaceID, kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicy, policyID), purge)
}

func (s *adminServer) listPoliciesByNamespace(ctx context.Context, namespaceID uuid.UUID) ([]*PolicyDoc, error) {
	partitionKey := azcosmos.NewPartitionKeyString(namespaceID.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.policyType FROM c
WHERE c.namespaceId = @namespaceId AND c.type = @type`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: namespaceID.String()},
				{Name: "@type", Value: kmsdoc.DocTypeNamePolicy},
			},
		})

	return PagerToList[PolicyDoc](ctx, pager)
}

func (doc *PolicyDoc) PopulatePolicyRef(r *PolicyRef) {
	r.ID = doc.GetUUID()
	r.NamespaceID = doc.NamespaceID
	r.Updated = doc.Updated
	r.UpdatedBy = doc.UpdatedBy
	r.Deleted = doc.Deleted

	r.PolicyType = doc.PolicyType
}

type PolicyStateStatus string

const (
	PolicyStateStatusSuccess PolicyStateStatus = "success"
)

type PolicyStateDoc struct {
	kmsdoc.BaseDoc
	PolicyType     PolicyType                           `json:"policyType"`
	Status         PolicyStateStatus                    `json:"status"`
	Message        string                               `json:"message"`
	CertRequest    *PolicyStateCertRequestDocSection    `json:"certRequest,omitempty"`
	CertAadAppCred *PolicyStateCertAadAppCredDocSection `json:"certAadAppCred,omitempty"`
}

func (s *adminServer) GetPolicyStateDoc(c context.Context, namespaceID uuid.UUID, policyID uuid.UUID) (*PolicyStateDoc, error) {
	pd := new(PolicyStateDoc)
	err := kmsdoc.AzCosmosRead(c, s.azCosmosContainerClientCerts, namespaceID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicyState, policyID), pd)
	return pd, err
}

func (doc *PolicyDoc) ToPolicy() *Policy {
	if doc == nil {
		return nil
	}
	p := Policy{
		ID:          doc.GetUUID(),
		PolicyType:  doc.PolicyType,
		NamespaceID: doc.NamespaceID,
		Deleted:     doc.Deleted,
		Updated:     doc.Updated,
		UpdatedBy:   fmt.Sprintf("%s:%s", doc.UpdatedBy, doc.UpdatedByName),
	}
	switch doc.PolicyType {
	case PolicyTypeCertRequest:
		p.CertRequest = doc.CertRequest.toCertificateRequestPolicyParameters()
	case PolicyTypeCertEnroll:
		p.CertEnroll = doc.CertEnroll.toCertificateEnrollPolicyParameters()
	case PolicyTypeCertAadAppClientCredential:
		p.CertAadAppCred = doc.CertAadAppCred.toCertificateAadAppPolicyParameters()
	}
	return &p
}

func (doc *PolicyStateDoc) ToPolicyState() *PolicyState {
	if doc == nil {
		return nil
	}
	ps := PolicyState{
		ID:          doc.GetUUID(),
		PolicyType:  doc.PolicyType,
		NamespaceID: doc.NamespaceID,
		Updated:     doc.Updated,
		UpdatedBy:   fmt.Sprintf("%s:%s", doc.UpdatedBy, doc.UpdatedByName),
	}
	switch doc.PolicyType {
	case PolicyTypeCertRequest:
		ps.CertRequest = doc.CertRequest.ToPolicyStateCertRequest()
	}
	return &ps
}
