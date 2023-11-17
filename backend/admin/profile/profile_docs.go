package profile

import (
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type ProfileDoc struct {
	resdoc.ResourceDoc

	DisplayName string `json:"displayName"`
}

type AppDoc struct {
	ProfileDoc

	ApplicationID        string `json:"applicationId,omitempty"`
	ServicePrincipalID   string `json:"servicePrincipalId,omitempty"`
	ServicePrincipalType string `json:"servicePrincipalType,omitempty"`
}

const (
	NamespaceIDApp   = "app"
	NamespaceIDCA    = "ca"
	NamespaceIDGraph = "graph"
)

func (doc *ProfileDoc) ToRef() (ref models.Ref) {
	ref = doc.ResourceDoc.ToRef()
	ref.DisplayName = doc.DisplayName
	return ref
}

func (doc *AppDoc) ToApplicationByAppId() (m models.ApplicationByAppId) {
	m.Ref = doc.ToRef()
	m.ApplicationId = doc.ApplicationID
	m.ServicePrincipalId = doc.ServicePrincipalID
	return m
}
