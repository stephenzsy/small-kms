package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

type contextKey string

const (
	appAuthIdentityContextKey contextKey = "smallkms.appAuthIdentity"
	roleKeyAppAdmin           string     = "App.Admin"
	roleKeyAgentActiveHost    string     = "Agent.ActiveHost"
)

type AuthIdentity interface {
	HasAppRole(role string) bool
	HasAdminRole() bool
	ClientPrincipalID() uuid.UUID
	ClientPrincipalName() string
	AppIDClaim() uuid.UUID
	GetOnBehalfOfTokenCredential(common.AzureAppConfidentialIdentity, *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error)
}

type authIdentity struct {
	msClientPrincipalID   uuid.UUID
	msClientPrincipalName string
	appRoles              map[string]bool
	appIDClaim            uuid.UUID
	bearerToken           string
}

func GetAuthIdentity(ctx context.Context) (identity AuthIdentity, ok bool) {
	ctxValue := ctx.Value(appAuthIdentityContextKey)
	identity, ok = ctxValue.(AuthIdentity)
	return
}

func (i authIdentity) HasAppRole(role string) bool {
	return i.appRoles[role]
}

func (i authIdentity) HasAdminRole() bool {
	return i.appRoles[roleKeyAppAdmin]
}

func (i authIdentity) ClientPrincipalID() uuid.UUID {
	return i.msClientPrincipalID
}

func (i authIdentity) ClientPrincipalName() string {
	return i.msClientPrincipalName
}

func (i authIdentity) AppIDClaim() uuid.UUID {
	return i.appIDClaim
}

func (i authIdentity) GetOnBehalfOfTokenCredential(s common.AzureAppConfidentialIdentity, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error) {
	return s.GetOnBehalfOfTokenCredential(i.bearerToken, opts)
}
