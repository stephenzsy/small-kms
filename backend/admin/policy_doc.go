package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyDoc struct {
	kmsdoc.BaseDoc
	PolicyType  PolicyType                   `json:"policyType"`
	CertRequest *PolicyCertRequestDocSection `json:"certRequest,omitempty"`
}

func (s *adminServer) GetPolicyDoc(c context.Context, namespaceID uuid.UUID, policyID uuid.UUID) (*PolicyDoc, error) {
	pd := new(PolicyDoc)
	err := kmsdoc.AzCosmosRead(c, s.azCosmosContainerClientCerts, namespaceID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicy, policyID), pd)
	return pd, err
}

type PolicyStateStatus string

const (
	PolicyStateStatusSuccess PolicyStateStatus = "success"
)

type PolicyStateDoc struct {
	kmsdoc.BaseDoc
	PolicyType  PolicyType                        `json:"policyType"`
	Status      PolicyStateStatus                 `json:"status"`
	Message     string                            `json:"message"`
	CertRequest *PolicyStateCertRequestDocSection `json:"certRequest,omitempty"`
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
		Updated:     doc.Updated,
		UpdatedBy:   fmt.Sprintf("%s:%s", doc.UpdatedBy, doc.UpdatedByName),
	}
	switch doc.PolicyType {
	case PolicyTypeCertRequest:
		p.CertRequest = doc.CertRequest.ToCertificateRequestPolicyParameters()
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
