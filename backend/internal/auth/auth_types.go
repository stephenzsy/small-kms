package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type contextKey int

const (
	authIdentityContextKey contextKey = iota
)

const (
	roleKeyAppAdmin string = "App.Admin"
)

type AuthIdentity interface {
	ClientPrincipalDisplayName() string
	HasAdminRole() bool
}

type authIdentity struct {
	msClientPrincipalID   uuid.UUID
	msClientPrincipalName string
	appRoles              map[string]bool
	bearerToken           string
}

// HasAdminRole implements AuthIdentity.
func (a *authIdentity) HasAdminRole() bool {
	return a.appRoles[roleKeyAppAdmin]
}

func GetAuthIdentity(c context.Context) AuthIdentity {
	return c.Value(authIdentityContextKey).(AuthIdentity)
}

func (a *authIdentity) ClientPrincipalDisplayName() string {
	return fmt.Sprintf("%s:%s", a.msClientPrincipalID, a.msClientPrincipalName)
}

var _ AuthIdentity = (*authIdentity)(nil)
