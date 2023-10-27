package profile

import (
	"github.com/google/uuid"
)

type GraphProfileDoc struct {
	ProfileDoc
}

type ServicePrincipalProfileDoc struct {
	GraphProfileDoc

	AppID uuid.UUID `json:"appId"`
}
