package ns

import (
	"fmt"
	"strings"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var (
	ErrInvalidNamespaceID = fmt.Errorf("invalid namespace id")
)

type NamespaceCertificateTemplateCapabilities struct {
	AllowedReservedNames       map[shared.Identifier]int
	AllowedIssuerNamespaces    utils.Set[shared.NamespaceIdentifier]
	DefaultIssuerNamespace     *shared.NamespaceIdentifier
	AllowedUsages              utils.Set[shared.CertificateUsage]
	AllowVariables             bool
	DefaultMaxValidityInMonths int
	DefaultKeyType             shared.JwtKty
	DefaultKeySize             int32
	DefaultRsaAlgorithm        shared.JwkAlg
	DefaultCrv                 shared.JwkCrv
	HasKeyStore                bool
	KeyExportable              bool
	RestrictKeyTypeRsa         bool
	DelegateForMembers         bool
}

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
	case shared.NamespaceKindSystem:
		return map[shared.Identifier]int{
			shared.StringIdentifier(shared.CertTemplateNameDefaultMtls): 0,
		}
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

var (
	caIntMsEntraNamespaceIdentifier = shared.NewNamespaceIdentifier(shared.NamespaceKindCaInt, shared.StringIdentifier(IntCaNameMsEntraClientSecret))
	caIntServiceNamespaceIdentifier = shared.NewNamespaceIdentifier(shared.NamespaceKindCaInt, shared.StringIdentifier(IntCaNameServices))
)

func GetAllowedCertificateIssuersForTemplate(templateLocator shared.ResourceLocator) (cap NamespaceCertificateTemplateCapabilities) {
	nsID := templateLocator.GetNamespaceID()
	templateID := templateLocator.GetID().Identifier()
	allowedNs := utils.NewSet[shared.NamespaceIdentifier]()
	allowedUsages := utils.NewSet[shared.CertificateUsage]()
	cap.DefaultMaxValidityInMonths = 12
	cap.DefaultKeyType = shared.KeyTypeRSA
	cap.DefaultKeySize = 2048
	cap.DefaultRsaAlgorithm = shared.AlgRS384
	cap.DefaultCrv = shared.CurveNameP384
	switch nsID.Kind() {
	case shared.NamespaceKindSystem:
		allowedNs.Add(nsID)
		allowedUsages.Add(shared.CertUsageClientAuth)
		cap.KeyExportable = false
	case shared.NamespaceKindCaRoot:
		allowedNs.Add(nsID)
		allowedUsages.Add(shared.CertUsageCA)
		allowedUsages.Add(shared.CertUsageCARoot)
		if nsID.Identifier().String() == string(RootCANameTest) {
			cap.DefaultMaxValidityInMonths = 6
			cap.DefaultKeyType = shared.KeyTypeEC
		} else {
			cap.DefaultMaxValidityInMonths = 120
			cap.DefaultKeySize = 4096
		}
		cap.HasKeyStore = true
		cap.KeyExportable = false
	case shared.NamespaceKindCaInt:
		if nsID.Identifier().String() == string(IntCaNameTest) {
			allowedNs.Add(shared.NewNamespaceIdentifier(shared.NamespaceKindCaRoot, shared.StringIdentifier(RootCANameTest)))
			cap.DefaultMaxValidityInMonths = 3
			cap.DefaultKeyType = shared.KeyTypeEC
		} else {
			allowedNs.Add(shared.NewNamespaceIdentifier(shared.NamespaceKindCaRoot, shared.StringIdentifier(RootCANameDefault)))
			cap.DefaultMaxValidityInMonths = 36
			cap.DefaultKeySize = 4096
		}
		cap.HasKeyStore = true
		cap.KeyExportable = false
		allowedUsages.Add(shared.CertUsageCA)
	case shared.NamespaceKindGroup:
		if strings.HasPrefix(templateID.String(), "test") {
			allowedNs.Add(shared.NewNamespaceIdentifier(shared.NamespaceKindCaInt, shared.StringIdentifier(IntCaNameTest)))
		}
		switch templateID.String() {
		case string(shared.CertTemplateNameDefaultIntranetAccess):
			allowedNs.Add(shared.NewNamespaceIdentifier(shared.NamespaceKindCaInt, shared.StringIdentifier(IntCaNameIntranet)))
			allowedUsages.Add(shared.CertUsageClientAuth)
			cap.DefaultMaxValidityInMonths = 1
			cap.HasKeyStore = false
		case string(shared.CertTemplateNameDefaultMsEntraClientCreds):
			allowedNs.Add(caIntMsEntraNamespaceIdentifier)
			allowedNs.Add(nsID)
			cap.DefaultIssuerNamespace = &caIntMsEntraNamespaceIdentifier
			allowedUsages.Add(shared.CertUsageClientAuth)
			allowedUsages.Add(shared.CertUsageServerAuth)
			cap.HasKeyStore = false
			cap.RestrictKeyTypeRsa = true
			cap.DefaultRsaAlgorithm = shared.AlgRS256
		}
		cap.AllowVariables = true
		cap.DelegateForMembers = true
	case shared.NamespaceKindServicePrincipal:
		if strings.HasPrefix(templateID.String(), "test") {
			allowedNs.Add(shared.NewNamespaceIdentifier(shared.NamespaceKindCaInt, shared.StringIdentifier(IntCaNameTest)))
		}
		switch templateID.String() {
		case string(shared.CertTemplateNameDefaultMsEntraClientCreds):
			allowedNs.Add(caIntMsEntraNamespaceIdentifier)
			allowedNs.Add(nsID)
			cap.DefaultIssuerNamespace = &caIntMsEntraNamespaceIdentifier
			cap.RestrictKeyTypeRsa = true
			cap.DefaultRsaAlgorithm = shared.AlgRS256
		case string(shared.CertTemplateNameDefaultMtls):
			allowedNs.Add(caIntServiceNamespaceIdentifier)
			allowedNs.Add(nsID)
			cap.DefaultIssuerNamespace = &caIntServiceNamespaceIdentifier
			cap.RestrictKeyTypeRsa = true
			cap.DefaultRsaAlgorithm = shared.AlgRS256
		default:
			allowedNs.Add(caIntServiceNamespaceIdentifier)
		}
		allowedUsages.Add(shared.CertUsageClientAuth)
		allowedUsages.Add(shared.CertUsageServerAuth)
		cap.HasKeyStore = true
		cap.KeyExportable = true
	}
	cap.AllowedReservedNames = GetReservedCertificateTemplateNames(nsID)
	cap.AllowedIssuerNamespaces = allowedNs
	cap.AllowedUsages = allowedUsages
	return
}

func validateNamespaceID(nsID shared.NamespaceIdentifier) error {
	switch nsID.Kind() {
	case shared.NamespaceKindSystem:
		switch nsID.Identifier().String() {
		case string(SystemServiceNameAgentPush):
			return nil
		}
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
type RequestContext = common.RequestContext

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
	return common.RequestContextWithValue(parent, namespaceContextKey, &namespaceContext{
		nsID: nsID,
	}), nil
}

func GetNamespaceContext(c RequestContext) NamespaceContext {
	if nc, ok := c.Value(namespaceContextKey).(NamespaceContext); ok {
		return nc
	}
	return nil
}
