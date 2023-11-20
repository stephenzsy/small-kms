package profile

import (
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type ProfileDoc struct {
	resdoc.ResourceDoc

	DisplayName          *string `json:"displayName,omitempty"`
	UserPrincipalName    *string `json:"userPrincipalName,omitempty"`
	ServicePrincipalType *string `json:"servicePrincipalType,omitempty"`
	AppId                *string `json:"appId,omitempty"`
	Mail                 *string `json:"mail,omitempty"`
}

type AppDoc struct {
	resdoc.ResourceDoc

	DisplayName          *string `json:"displayName,omitempty"`
	ApplicationID        *string `json:"applicationId,omitempty"`
	ServicePrincipalID   *string `json:"servicePrincipalId,omitempty"`
	ServicePrincipalType *string `json:"servicePrincipalType,omitempty"`
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

func (doc *AppDoc) ToProfile() (m models.Profile) {
	m.Ref = doc.ToRef()
	m.DisplayName = doc.DisplayName
	m.ApplicationId = doc.ApplicationID
	m.ServicePrincipalId = doc.ServicePrincipalID
	return m
}

func (doc *ProfileDoc) ToModel() (m models.Profile) {
	m.Ref = doc.ToRef()
	m.UserPrincipalName = doc.UserPrincipalName
	m.ServicePrincipalType = doc.ServicePrincipalType
	m.AppId = doc.AppId
	m.Mail = doc.Mail
	return m
}

func (doc *ProfileDoc) TargetNamespaceProvider() models.NamespaceProvider {
	switch doc.PartitionKey.ResourceProvider {
	case models.ProfileResourceProviderServicePrincipal:
		return models.NamespaceProviderServicePrincipal
	case models.ProfileResourceProviderGroup:
		return models.NamespaceProviderGroup
	case models.ProfileResourceProviderUser:
		return models.NamespaceProviderUser
	}
	return ""
}
