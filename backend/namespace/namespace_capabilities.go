package ns

import (
	"fmt"
	"strings"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var (
	ErrInvalidNamespaceID = fmt.Errorf("invalid namespace id")
)

type NamespaceID = kmsdoc.DocNsID

type NamespaceCertificateTemplateCapabilities struct {
	AllowedReservedNames       map[common.Identifier]int
	AllowedIssuerNamespaces    utils.Set[NamespaceID]
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

type NamespaceCapabilities interface {
	GetAllowedCertificateIssuersForTemplate(templateID common.Identifier, expectedProfileType models.ProfileType) NamespaceCertificateTemplateCapabilities
	GetReservedCertificateTemplateNames(expectedProfileType models.ProfileType) map[common.Identifier]int
}

type namespaceCapabilities struct {
	nsID NamespaceID // must be validated
}

func (nc *namespaceCapabilities) GetReservedCertificateTemplateNames(expectedProfileType models.ProfileType) (r map[common.Identifier]int) {
	switch nc.nsID.Kind() {
	case kmsdoc.DocNsTypeCaRoot,
		kmsdoc.DocNsTypeCaInt:
		return map[common.Identifier]int{
			common.StringIdentifier(string(CertTemplateNameDefault)): 0,
		}
	case kmsdoc.DocNSTypeDirectory:
		switch expectedProfileType {
		case models.ProfileTypeGroup:
			return map[common.Identifier]int{
				common.StringIdentifier(string(CertTemplateNameDefaultMsEntraClientCreds)): 0,
				common.StringIdentifier(string(CertTemplateNameDefaultIntranetAccess)):     1,
			}
		case models.ProfileTypeServicePrincipal:
			return map[common.Identifier]int{
				common.StringIdentifier(string(CertTemplateNameDefault)):                   0,
				common.StringIdentifier(string(CertTemplateNameDefaultMsEntraClientCreds)): 1,
			}
		}
	}
	return
}

func (nc *namespaceCapabilities) GetAllowedCertificateIssuersForTemplate(templateID common.Identifier, expectedProfileType models.ProfileType) (cap NamespaceCertificateTemplateCapabilities) {
	allowedNs := utils.NewSet[NamespaceID]()
	allowedUsages := utils.NewSet[models.CertificateUsage]()
	cap.DefaultMaxValidityInMonths = 12
	cap.DefaultKeyType = models.KeyTypeRSA
	cap.DefaultKeySize = 2048
	cap.DefaultRsaAlgorithm = models.AlgRS384
	cap.DefaultCrv = models.CurveNameP384
	switch nc.nsID.Kind() {
	case kmsdoc.DocNsTypeCaRoot:
		allowedNs.Add(nc.nsID)
		allowedUsages.Add(models.CertUsageCA)
		allowedUsages.Add(models.CertUsageCARoot)
		cap.SelfSigned = true
		if nc.nsID.Identifier().String() == string(RootCANameTest) {
			cap.DefaultMaxValidityInMonths = 6
			cap.DefaultKeyType = models.KeyTypeEC
		} else {
			cap.DefaultMaxValidityInMonths = 120
			cap.DefaultKeySize = 4096
		}
		cap.HasKeyStore = true
		cap.KeyExportable = false
	case kmsdoc.DocNsTypeCaInt:
		if nc.nsID.Identifier().String() == string(IntCaNameTest) {
			allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaRoot, common.StringIdentifier(string(RootCANameTest))))
			cap.DefaultMaxValidityInMonths = 3
			cap.DefaultKeyType = models.KeyTypeEC
		} else {
			allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaRoot, common.StringIdentifier(string(RootCANameDefault))))
			cap.DefaultMaxValidityInMonths = 36
			cap.DefaultKeySize = 4096
		}
		cap.HasKeyStore = true
		cap.KeyExportable = false
		allowedUsages.Add(models.CertUsageCA)
	case kmsdoc.DocNSTypeDirectory:
		switch expectedProfileType {
		case models.ProfileTypeGroup:
			if strings.HasPrefix(templateID.String(), "test") {
				allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, common.StringIdentifier(string(IntCaNameTest))))
			}
			switch templateID.String() {
			case string(CertTemplateNameDefaultIntranetAccess):
				allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, common.StringIdentifier(string(IntCaNameIntranet))))
				allowedUsages.Add(models.CertUsageClientAuth)
				cap.DefaultMaxValidityInMonths = 1
				cap.HasKeyStore = false
			case string(CertTemplateNameDefaultMsEntraClientCreds):
				allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, common.StringIdentifier(string(IntCaNameMsEntraClientSecret))))
				allowedUsages.Add(models.CertUsageClientAuth)
				allowedUsages.Add(models.CertUsageServerAuth)
				cap.HasKeyStore = false
				cap.RestrictKeyTypeRsa = true
				cap.DefaultRsaAlgorithm = models.AlgRS256
			}
			cap.AllowVariables = true
			cap.DelegateForMembers = true
		case models.ProfileTypeServicePrincipal:
			if strings.HasPrefix(templateID.String(), "test") {
				allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, common.StringIdentifier(string(IntCaNameTest))))
			}
			switch templateID.String() {
			case string(CertTemplateNameDefaultMsEntraClientCreds):
				allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, common.StringIdentifier(string(IntCaNameMsEntraClientSecret))))
				cap.RestrictKeyTypeRsa = true
				cap.DefaultRsaAlgorithm = models.AlgRS256
			default:
				allowedNs.Add(kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, common.StringIdentifier(string(IntCaNameServices))))
			}
			allowedUsages.Add(models.CertUsageClientAuth)
			allowedUsages.Add(models.CertUsageServerAuth)
			cap.HasKeyStore = true
			cap.KeyExportable = true
		}
	}
	cap.AllowedReservedNames = nc.GetReservedCertificateTemplateNames(expectedProfileType)
	cap.AllowedIssuerNamespaces = allowedNs
	cap.AllowedUsages = allowedUsages
	return
}

func validateNamespaceID(nsID NamespaceID) error {
	switch nsID.Kind() {
	case kmsdoc.DocNsTypeCaRoot:
		switch nsID.Identifier().String() {
		case string(RootCANameDefault),
			string(RootCANameTest):
			return nil
		}
	case kmsdoc.DocNsTypeCaInt:
		switch nsID.Identifier().String() {
		case string(IntCaNameServices),
			string(IntCaNameIntranet),
			string(IntCaNameMsEntraClientSecret),
			string(IntCaNameTest):
			return nil
		}
	case kmsdoc.DocNSTypeDirectory:
		if id, isUuid := nsID.Identifier().TryGetUUID(); isUuid && id.Version() == 4 {
			return nil
		}
	case kmsdoc.DocNsTypeProfile:
		switch nsID.Identifier().String() {
		case string(ProfileNamespaceIDNameBuiltin),
			string(ProfileNamespaceIDNameTenant):
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrInvalidNamespaceID, nsID.String())
}

func GetNamespaceCapabilities(nsID NamespaceID) (NamespaceCapabilities, error) {
	if err := validateNamespaceID(nsID); err != nil {
		return nil, err
	}
	return &namespaceCapabilities{
		nsID: nsID,
	}, nil
}
