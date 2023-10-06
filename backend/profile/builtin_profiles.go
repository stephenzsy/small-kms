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
	idIntCaServices            = common.StringIdentifier("services")
	idIntCaIntranet            = common.StringIdentifier("intranet")
	idIntCaMsEntraClientSecret = common.StringIdentifier("ms-entra-client-secret")
	idIntCaTest                = common.StringIdentifier("test")
)

var (
	docNsIDProfileBuiltIn kmsdoc.DocNsID = kmsdoc.StringDocIdentifier(kmsdoc.DocNsTypeProfile, "builtin")
	docNsIDProfileTenant  kmsdoc.DocNsID = kmsdoc.StringDocIdentifier(kmsdoc.DocNsTypeProfile, "tenant")
)

var rootCaProfileDocs = map[common.Identifier]ProfileDoc{
	idCaRoot: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          kmsdoc.NewDocIdentifier(kmsdoc.DocKindCaRoot, idCaRoot),
		},
		DispalyName: utils.ToPtr("Default Root CA"),
		ProfileType: models.ProfileTypeRootCA,
	},
	idCaRootTest: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          kmsdoc.NewDocIdentifier(kmsdoc.DocKindCaRoot, idCaRootTest),
		},
		DispalyName: utils.ToPtr("Test Root CA"),
		ProfileType: models.ProfileTypeRootCA,
	},
}

var intCaProfileDocs = map[common.Identifier]ProfileDoc{
	idIntCaServices: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          kmsdoc.NewDocIdentifier(kmsdoc.DocKindCaInt, idIntCaServices),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Services"),
		ProfileType: models.ProfileTypeIntermediateCA,
	},
	idIntCaIntranet: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          kmsdoc.NewDocIdentifier(kmsdoc.DocKindCaInt, idIntCaIntranet),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Intranet Access"),
		ProfileType: models.ProfileTypeIntermediateCA,
	},
	idIntCaMsEntraClientSecret: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          kmsdoc.NewDocIdentifier(kmsdoc.DocKindCaInt, idIntCaMsEntraClientSecret),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Microsoft Entra Client Secert"),
		ProfileType: models.ProfileTypeIntermediateCA,
	},
	idIntCaTest: {
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: docNsIDProfileBuiltIn,
			ID:          kmsdoc.NewDocIdentifier(kmsdoc.DocKindCaInt, idIntCaTest),
		},
		DispalyName: utils.ToPtr("Intermediate CA - Test"),
		ProfileType: models.ProfileTypeIntermediateCA,
	},
}
