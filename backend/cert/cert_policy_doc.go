package cert

import (
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/key"
)

type CertPolicyDoc struct {
	base.BaseDoc

	DisplayName     string                   `json:"displayName"`
	KeySpec         key.SigningKeySpec       `json:"keySpec"`
	KeyExportable   bool                     `json:"keyExportable"`
	ExpiryTime      base.Period              `json:"expiryTime"`
	LifetimeAction  *key.LifetimeAction      `json:"lifetimeActions,omitempty"`
	Subject         CertificateSubject       `json:"subject"`
	SANs            *SubjectAlternativeNames `json:"sans,omitempty"`
	Flags           []CertificateFlag        `json:"flags"`
	Version         HexDigest                `json:"version"`
	IssuerNamespace base.NamespaceIdentifier `json:"issuerNamespace"`
}

const (
	queryColumnDisplayName = "c.displayName"
)

// populate ref
func (d *CertPolicyDoc) PopulateModelRef(m *CertPolicyRef) {
	if d == nil || m == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&m.ResourceReference)
	m.DisplayName = d.DisplayName
}

func (d *CertPolicyDoc) PopulateModel(m *CertPolicy) {
	if d == nil || m == nil {
		return
	}
	d.PopulateModelRef(&m.CertPolicyRef)
	m.KeySpec = d.KeySpec
	if m.KeySpec.KeyOperations == nil {
		m.KeySpec.KeyOperations = []key.JsonWebKeyOperation{}
	}
	m.KeyExportable = d.KeyExportable
	m.ExpiryTime = d.ExpiryTime
	m.LifetimeAction = d.LifetimeAction
	m.Subject = d.Subject
	m.SubjectAlternativeNames = d.SANs
	m.Flags = d.Flags
	m.Version = d.Version
	m.IssuerNamespaceKind = d.IssuerNamespace.Kind()
	m.IssuerNamespaceIdentifier = d.IssuerNamespace.ID()
}

var _ base.ModelRefPopulater[CertPolicyRef] = (*CertPolicyDoc)(nil)
var _ base.ModelPopulater[CertPolicy] = (*CertPolicyDoc)(nil)
