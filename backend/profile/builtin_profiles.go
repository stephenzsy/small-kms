package profile

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var (
	idCaRoot     = common.StringIdentifier(string(ns.RootCANameDefault))
	idCaRootTest = common.StringIdentifier(string(ns.RootCANameTest))
)

var (
	idIntCaServices            = common.StringIdentifier(ns.IntCaNameServices)
	idIntCaIntranet            = common.StringIdentifier(ns.IntCaNameIntranet)
	idIntCaMsEntraClientSecret = common.StringIdentifier(ns.IntCaNameMsEntraClientSecret)
	idIntCaTest                = common.StringIdentifier(ns.IntCaNameTest)
)

var (
	docNsIDProfileBuiltIn models.NamespaceID = common.NewIdentifierWithKind(models.NamespaceKindProfile, common.StringIdentifier(ns.ProfileNamespaceIDNameBuiltin))
	docNsIDProfileTenant  models.NamespaceID = common.NewIdentifierWithKind(models.NamespaceKindProfile, common.StringIdentifier(ns.ProfileNamespaceIDNameTenant))
)

var rootCaProfileDocs = map[common.Identifier]ProfileDoc{
	idCaRoot: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCaRoot, idCaRoot),
		},
		DispalyName: utils.ToPtr("Default Root CA"),
		ProfileType: models.NamespaceKindCaRoot,
	},
	idCaRootTest: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCaRoot, idCaRootTest),
		},
		DispalyName: utils.ToPtr("Test Root CA"),
		ProfileType: models.NamespaceKindCaRoot,
	},
}

var intCaProfileDocs = map[common.Identifier]ProfileDoc{
	idIntCaServices: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCaInt, idIntCaServices),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Services"),
		ProfileType: models.NamespaceKindCaInt,
	},
	idIntCaIntranet: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCaInt, idIntCaIntranet),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Intranet Access"),
		ProfileType: models.NamespaceKindCaInt,
	},
	idIntCaMsEntraClientSecret: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCaInt, idIntCaMsEntraClientSecret),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Microsoft Entra Client Secert"),
		ProfileType: models.NamespaceKindCaInt,
	},
	idIntCaTest: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          common.NewIdentifierWithKind(models.ResourceKindCaInt, idIntCaTest),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Test"),
		ProfileType: models.NamespaceKindCaInt,
	},
}
