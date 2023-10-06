package profile

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

var (
	idCaRoot     = models.IdentifierFromString("default")
	idCaRootTest = models.IdentifierFromString("test")
)

var (
	idIntCaServices            = models.IdentifierFromString("services")
	idIntCaIntranet            = models.IdentifierFromString("intranet")
	idIntCaMsEntraClientSecret = models.IdentifierFromString("ms-entra-client-secret")
	idIntCaTest                = models.IdentifierFromString("test")
)

var rootCaProfiles = map[models.Identifier]models.Profile{
	idCaRoot: {
		Type:        models.ProfileTypeRootCA,
		Identifier:  idCaRoot,
		DisplayName: "Default Root CA",
	},
	idCaRootTest: {
		Type:        models.ProfileTypeRootCA,
		Identifier:  idCaRootTest,
		DisplayName: "Test Root CA",
	},
}

var intCaProfiles = map[models.Identifier]models.Profile{
	idIntCaServices: {
		Type:        models.ProfileTypeIntermediateCA,
		Identifier:  idIntCaServices,
		DisplayName: "Intermediate CA - Services",
	},
	idIntCaIntranet: {
		Type:        models.ProfileTypeIntermediateCA,
		Identifier:  idIntCaIntranet,
		DisplayName: "Intermediate CA - Intranet Access",
	},
	idIntCaMsEntraClientSecret: {
		Type:        models.ProfileTypeIntermediateCA,
		Identifier:  idIntCaMsEntraClientSecret,
		DisplayName: "Intermediate CA - Microsoft Entra Client Secert",
	},
	idIntCaTest: {
		Type:        models.ProfileTypeIntermediateCA,
		Identifier:  idIntCaTest,
		DisplayName: "Intermediate CA - Test",
	},
}

func getBuiltInCaProfiles() []models.ProfileRef {
	return []models.ProfileRef{
		rootCaProfiles[idCaRoot],
		rootCaProfiles[idCaRootTest],
	}
}

func getBuiltInIntermediateCaProfiles() []models.ProfileRef {
	return []models.ProfileRef{
		intCaProfiles[idIntCaServices],
		intCaProfiles[idIntCaIntranet],
		intCaProfiles[idIntCaMsEntraClientSecret],
		intCaProfiles[idIntCaTest],
	}
}

// ListProfiles implements ProfileService.
func (*profileService) ListProfiles(c common.ServiceContext, profileType models.ProfileType) ([]models.ProfileRef, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	switch profileType {
	case models.ProfileTypeRootCA:
		return getBuiltInCaProfiles(), nil
	case models.ProfileTypeIntermediateCA:
		return getBuiltInIntermediateCaProfiles(), nil
	}
	return make([]models.ProfileRef, 0), nil
	//	panic("unimplemented")
}
