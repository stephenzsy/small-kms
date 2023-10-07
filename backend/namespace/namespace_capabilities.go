package ns

import (
	"context"
	"fmt"
	"strings"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var (
	ErrInvalidNamespaceID = fmt.Errorf("invalid namespace id")
)

type NamespaceCertificateTemplateCapabilities struct {
	AllowedReservedNames       map[common.Identifier]int
	AllowedIssuerNamespaces    utils.Set[models.NamespaceID]
	AllowedUsages              utils.Set[models.CertificateUsage]
	AllowVariables             bool
	SelfSigned                 bool
	DefaultMaxValidityInMonths int
	DefaultKeyType             models.JwtKty
	DefaultKeySize             int
	DefaultRsaAlgorithm        models.JwkAlg
	DefaultCrv                 models.JwtCrv
	HasKeyStore                bool
	KeyExportable              bool
	RestrictKeyTypeRsa         bool
	DelegateForMembers         bool
}

type NamespaceContext interface {
	GetID() models.NamespaceID
}

type namespaceContext struct {
	nsID models.NamespaceID // must be validated
}

func (nc *namespaceContext) GetID() models.NamespaceID {
	return nc.nsID
}

func GetReservedCertificateTemplateNames(nsID models.NamespaceID) (r map[common.Identifier]int) {
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot,
		models.NamespaceKindCaInt:
		return map[common.Identifier]int{
			common.StringIdentifier(string(CertTemplateNameDefault)): 0,
		}
	case models.NamespaceKindGroup:
		return map[common.Identifier]int{
			common.StringIdentifier(string(CertTemplateNameDefaultMsEntraClientCreds)): 0,
			common.StringIdentifier(string(CertTemplateNameDefaultIntranetAccess)):     1,
		}
	case models.NamespaceKindServicePrincipal:
		return map[common.Identifier]int{
			common.StringIdentifier(string(CertTemplateNameDefault)):                   0,
			common.StringIdentifier(string(CertTemplateNameDefaultMsEntraClientCreds)): 1,
		}
	}
	return
}

func GetAllowedCertificateIssuersForTemplate(templateLocator models.ResourceLocator) (cap NamespaceCertificateTemplateCapabilities) {
	nsID := templateLocator.GetNamespaceID()
	templateID := templateLocator.GetID().Identifier()
	allowedNs := utils.NewSet[models.NamespaceID]()
	allowedUsages := utils.NewSet[models.CertificateUsage]()
	cap.DefaultMaxValidityInMonths = 12
	cap.DefaultKeyType = models.KeyTypeRSA
	cap.DefaultKeySize = 2048
	cap.DefaultRsaAlgorithm = models.AlgRS384
	cap.DefaultCrv = models.CurveNameP384
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot:
		allowedNs.Add(nsID)
		allowedUsages.Add(models.CertUsageCA)
		allowedUsages.Add(models.CertUsageCARoot)
		cap.SelfSigned = true
		if nsID.Identifier().String() == string(RootCANameTest) {
			cap.DefaultMaxValidityInMonths = 6
			cap.DefaultKeyType = models.KeyTypeEC
		} else {
			cap.DefaultMaxValidityInMonths = 120
			cap.DefaultKeySize = 4096
		}
		cap.HasKeyStore = true
		cap.KeyExportable = false
	case models.NamespaceKindCaInt:
		if nsID.Identifier().String() == string(IntCaNameTest) {
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaRoot, common.StringIdentifier(RootCANameTest)))
			cap.DefaultMaxValidityInMonths = 3
			cap.DefaultKeyType = models.KeyTypeEC
		} else {
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaRoot, common.StringIdentifier(RootCANameDefault)))
			cap.DefaultMaxValidityInMonths = 36
			cap.DefaultKeySize = 4096
		}
		cap.HasKeyStore = true
		cap.KeyExportable = false
		allowedUsages.Add(models.CertUsageCA)
	case models.NamespaceKindGroup:
		if strings.HasPrefix(templateID.String(), "test") {
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaInt, common.StringIdentifier(IntCaNameTest)))
		}
		switch templateID.String() {
		case string(CertTemplateNameDefaultIntranetAccess):
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaInt, common.StringIdentifier(IntCaNameIntranet)))
			allowedUsages.Add(models.CertUsageClientAuth)
			cap.DefaultMaxValidityInMonths = 1
			cap.HasKeyStore = false
		case string(CertTemplateNameDefaultMsEntraClientCreds):
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaInt, common.StringIdentifier(IntCaNameMsEntraClientSecret)))
			allowedUsages.Add(models.CertUsageClientAuth)
			allowedUsages.Add(models.CertUsageServerAuth)
			cap.HasKeyStore = false
			cap.RestrictKeyTypeRsa = true
			cap.DefaultRsaAlgorithm = models.AlgRS256
		}
		cap.AllowVariables = true
		cap.DelegateForMembers = true
	case models.NamespaceKindServicePrincipal:
		if strings.HasPrefix(templateID.String(), "test") {
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaInt, common.StringIdentifier(IntCaNameTest)))
		}
		switch templateID.String() {
		case string(CertTemplateNameDefaultMsEntraClientCreds):
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaInt, common.StringIdentifier(IntCaNameMsEntraClientSecret)))
			cap.RestrictKeyTypeRsa = true
			cap.DefaultRsaAlgorithm = models.AlgRS256
		default:
			allowedNs.Add(common.NewIdentifierWithKind(models.NamespaceKindCaInt, common.StringIdentifier(IntCaNameServices)))
		}
		allowedUsages.Add(models.CertUsageClientAuth)
		allowedUsages.Add(models.CertUsageServerAuth)
		cap.HasKeyStore = true
		cap.KeyExportable = true

	}
	cap.AllowedReservedNames = GetReservedCertificateTemplateNames(nsID)
	cap.AllowedIssuerNamespaces = allowedNs
	cap.AllowedUsages = allowedUsages
	return
}

func validateNamespaceID(nsID models.NamespaceID) error {
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot:
		switch nsID.Identifier().String() {
		case string(RootCANameDefault),
			string(RootCANameTest):
			return nil
		}
	case models.NamespaceKindCaInt:
		switch nsID.Identifier().String() {
		case string(IntCaNameServices),
			string(IntCaNameIntranet),
			string(IntCaNameMsEntraClientSecret),
			string(IntCaNameTest):
			return nil
		}
	case models.NamespaceKindGroup,
		models.NamespaceKindApplication,
		models.NamespaceKindDevice,
		models.NamespaceKindServicePrincipal,
		models.NamespaceKindUser:
		if id, isUuid := nsID.Identifier().TryGetUUID(); isUuid && id.Version() == 4 {
			return nil
		}
	case models.NamespaceKindProfile:
		switch nsID.Identifier().String() {
		case string(ProfileNamespaceIDNameBuiltin),
			string(ProfileNamespaceIDNameTenant):
			return nil
		}
	}
	return fmt.Errorf("%w:invalid namespace ID: %s", common.ErrStatusBadRequest, nsID.String())
}

type contextKey string

const (
	namespaceContextKey contextKey = "namespaceContext"
)

func WithNamespaceContext(parent common.ServiceContext, unverifiedKind models.NamespaceKind, unverifiedIdentifier common.Identifier) (common.ServiceContext, error) {
	if !unverifiedIdentifier.IsValid() {
		return parent, fmt.Errorf("%w:invalid namespace identifier", common.ErrStatusBadRequest)
	}
	nsID := common.NewIdentifierWithKind(unverifiedKind, unverifiedIdentifier)
	if err := validateNamespaceID(nsID); err != nil {
		return parent, err
	}
	return context.WithValue(parent, namespaceContextKey, &namespaceContext{
		nsID: nsID,
	}), nil
}

func GetNamespaceContext(c common.ServiceContext) NamespaceContext {
	if nc, ok := c.Value(namespaceContextKey).(NamespaceContext); ok {
		return nc
	}
	return nil
}
