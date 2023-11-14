package agentadmin

import "github.com/stephenzsy/small-kms/backend/resdoc"

type AgentDoc struct {
	resdoc.ResourceDoc

	DisplayName          string `json:"displayName"`
	ApplicationID        string `json:"applicationId,omitempty"`
	ServicePrincipalID   string `json:"servicePrincipalId,omitempty"`
	ServicePrincipalType string `json:"servicePrincipalType,omitempty"`
}
