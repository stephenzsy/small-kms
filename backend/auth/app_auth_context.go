package auth

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const (
	appAuthIdentityContextKey string = "smallkms.appAuthIdentity"
	roleKeyAppAdmin           string = "App.Admin"
)

type AuthIdentity interface {
	HasAdminRole() bool
	ClientPrincipalID() uuid.UUID
	ClientPrincipalName() string
	AppIDClaim() uuid.UUID
}

type authIdentity struct {
	msClientPrincipalID   uuid.UUID
	msClientPrincipalName string
	appRoles              map[string]bool
	appIDClaim            uuid.UUID
}

func GetAuthIdentity(ctx context.Context) (identity AuthIdentity, ok bool) {
	ctxValue := ctx.Value(appAuthIdentityContextKey)
	identity, ok = ctxValue.(AuthIdentity)
	return
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
