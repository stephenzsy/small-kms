package ns

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/shared"
)

var (
	ErrInvalidNamespaceID = fmt.Errorf("invalid namespace id")
)

type NamespaceContext interface {
	GetID() shared.NamespaceIdentifier
}

type namespaceContext struct {
	nsID shared.NamespaceIdentifier // must be validated
}

func (nc *namespaceContext) GetID() shared.NamespaceIdentifier {
	return nc.nsID
}

func GetReservedCertificateTemplateNames(nsID shared.NamespaceIdentifier) (r map[shared.Identifier]int) {
	switch nsID.Kind() {
	case shared.NamespaceKindCaRoot,
		shared.NamespaceKindCaInt:
		return map[shared.Identifier]int{
			shared.StringIdentifier(shared.CertTemplateNameDefault): 0,
		}
	case shared.NamespaceKindGroup:
		return map[shared.Identifier]int{
			shared.StringIdentifier(shared.CertTemplateNameDefaultMsEntraClientCreds): 0,
			shared.StringIdentifier(shared.CertTemplateNameDefaultIntranetAccess):     1,
		}
	case shared.NamespaceKindServicePrincipal:
		return map[shared.Identifier]int{
			shared.StringIdentifier(shared.CertTemplateNameDefault):                   0,
			shared.StringIdentifier(shared.CertTemplateNameDefaultMsEntraClientCreds): 1,
			shared.StringIdentifier(shared.CertTemplateNameDefaultMtls):               2,
		}
	}
	return
}

func validateNamespaceID(nsID shared.NamespaceIdentifier) error {
	switch nsID.Kind() {
	case shared.NamespaceKindCaRoot:
		switch nsID.Identifier().String() {
		case string(RootCANameDefault),
			string(RootCANameTest):
			return nil
		}
	case shared.NamespaceKindCaInt:
		switch nsID.Identifier().String() {
		case string(IntCaNameServices),
			string(IntCaNameIntranet),
			string(IntCaNameMsEntraClientSecret),
			string(IntCaNameTest):
			return nil
		}
	case shared.NamespaceKindGroup,
		shared.NamespaceKindApplication,
		shared.NamespaceKindDevice,
		shared.NamespaceKindServicePrincipal,
		shared.NamespaceKindUser:
		if id, isUuid := nsID.Identifier().TryGetUUID(); isUuid && id.Version() == 4 {
			return nil
		}
	case shared.NamespaceKindProfile:
		switch nsID.Identifier().String() {
		case string(ProfileNamespaceIDNameBuiltin),
			string(ProfileNamespaceIDNameTenant):
			return nil
		}
	}
	return fmt.Errorf("%w:invalid namespace ID: %s", common.ErrStatusBadRequest, nsID.String())
}

type contextKey string
type RequestContext = ctx.RequestContext

const (
	namespaceContextKey contextKey = "namespaceContext"
)

func WithNamespaceContext(parent RequestContext, unverifiedKind shared.NamespaceKind, unverifiedIdentifier shared.Identifier) (RequestContext, error) {
	if !unverifiedIdentifier.IsValid() {
		return parent, fmt.Errorf("%w:invalid namespace identifier", common.ErrStatusBadRequest)
	}
	nsID := shared.NewNamespaceIdentifier(unverifiedKind, unverifiedIdentifier)
	if err := validateNamespaceID(nsID); err != nil {
		return parent, err
	}
	return parent.WithValue(namespaceContextKey, &namespaceContext{
		nsID: nsID,
	}), nil
}

func GetNamespaceContext(c context.Context) NamespaceContext {
	if nc, ok := c.Value(namespaceContextKey).(NamespaceContext); ok {
		return nc
	}
	return nil
}
