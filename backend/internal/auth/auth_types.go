package auth

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
)

type contextKey int

const (
	authIdentityContextKey contextKey = iota
	AppConfidentialIdentityContextKey
)

const (
	roleValueAppAdmin        string = "App.Admin"
	RoleValueAgentActiveHost        = "Agent.ActiveHost"
)

type AuthIdentity interface {
	ClientPrincipalID() uuid.UUID
	ClientPrincipalDisplayName() string
	AppID() string
	HasAdminRole() bool
	HasRole(roleValue string) bool
	GetOnBehalfOfTokenCredential(c context.Context, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error)
}

type authIdentity struct {
	msClientPrincipalID   uuid.UUID
	msClientPrincipalName string
	appRoles              map[string]bool
	bearerToken           string
	appID                 string
}

// HasRole implements AuthIdentity.
func (a *authIdentity) HasRole(roleValue string) bool {
	if v, ok := a.appRoles[roleValue]; ok {
		return v
	}
	return false
}

// AppID implements AuthIdentity.
func (a *authIdentity) AppID() string {
	return a.appID
}

type AzureIdentity interface {
	TokenCredential() azcore.TokenCredential
	TenantID() string
	ClientID() string
}

type AzureAppConfidentialIdentity interface {
	AzureIdentity
	NewOnBehalfOfTokenCredential(userAssertion string, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error)
}

// ClientPrincipalID implements AuthIdentity.
func (a *authIdentity) ClientPrincipalID() uuid.UUID {
	return a.msClientPrincipalID
}

// GetOnBehalfOfTokenCredential implements AuthIdentity.
func (a *authIdentity) GetOnBehalfOfTokenCredential(c context.Context, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error) {
	i := c.Value(AppConfidentialIdentityContextKey).(AzureAppConfidentialIdentity)
	return i.NewOnBehalfOfTokenCredential(a.bearerToken, opts)
}

// HasAdminRole implements AuthIdentity.
func (a *authIdentity) HasAdminRole() bool {
	return a.appRoles[roleValueAppAdmin]
}

func GetAuthIdentity(c context.Context) AuthIdentity {
	return c.Value(authIdentityContextKey).(AuthIdentity)
}

func (a *authIdentity) ClientPrincipalDisplayName() string {
	return fmt.Sprintf("%s:%s", a.msClientPrincipalID, a.msClientPrincipalName)
}

var _ AuthIdentity = (*authIdentity)(nil)
