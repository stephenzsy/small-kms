package profile

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
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
	docNsIDProfileBuiltIn models.NamespaceID = shared.NewNamespaceIdentifier(shared.NamespaceKindProfile, common.StringIdentifier(ns.ProfileNamespaceIDNameBuiltin))
	docNsIDProfileTenant  models.NamespaceID = shared.NewNamespaceIdentifier(shared.NamespaceKindProfile, common.StringIdentifier(ns.ProfileNamespaceIDNameTenant))
)

var rootCaProfileDocs = map[common.Identifier]ProfileDoc{
	idCaRoot: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          shared.NewResourceIdentifier(shared.ResourceKindCaRoot, idCaRoot),
		},
		DispalyName: "Default Root CA",
		ProfileType: shared.NamespaceKindCaRoot,
		IsBuiltIn:   true,
	},
	idCaRootTest: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          shared.NewResourceIdentifier(shared.ResourceKindCaRoot, idCaRootTest),
		},
		DispalyName: "Test Root CA",
		ProfileType: shared.NamespaceKindCaRoot,
		IsBuiltIn:   true,
	},
}

var intCaProfileDocs = map[common.Identifier]ProfileDoc{
	idIntCaServices: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          shared.NewResourceIdentifier(shared.ResourceKindCaInt, idIntCaServices),
		},
		DispalyName: "Intermediate CA - Services",
		ProfileType: shared.NamespaceKindCaInt,
		IsBuiltIn:   true,
	},
	idIntCaIntranet: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          shared.NewResourceIdentifier(shared.ResourceKindCaInt, idIntCaIntranet),
		},
		DispalyName: "Intermediate CA - Intranet Access",
		ProfileType: shared.NamespaceKindCaInt,
		IsBuiltIn:   true,
	},
	idIntCaMsEntraClientSecret: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          shared.NewResourceIdentifier(shared.ResourceKindCaInt, idIntCaMsEntraClientSecret),
		},
		DispalyName: "Intermediate CA - Microsoft Entra Client Secert",
		ProfileType: shared.NamespaceKindCaInt,
		IsBuiltIn:   true,
	},
	idIntCaTest: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          shared.NewResourceIdentifier(shared.ResourceKindCaInt, idIntCaTest),
		},
		DispalyName: "Intermediate CA - Test",
		ProfileType: shared.NamespaceKindCaInt,
		IsBuiltIn:   true,
	},
}
