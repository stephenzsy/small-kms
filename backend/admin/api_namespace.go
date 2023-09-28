package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func getRootCaRefs() []Ref {
	return []Ref{
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_RootCA, DisplayName: "Root CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestRootCA, DisplayName: "Test Root CA", Type: RefTypeNamespace},
	}
}

func getIntCaRefs() []Ref {
	return []Ref{
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAService, DisplayName: "Services Intermediate CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_IntCAIntranet, DisplayName: "Intranet Intermediate CA", Type: RefTypeNamespace},
		{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestIntCA, DisplayName: "Test Intermediate CA", Type: RefTypeNamespace},
	}
}

/*
	func getBuiltInCaRootNamespaceRefs() []NamespaceRef {
		return []NamespaceRef{
			{NamespaceID: uuid.Nil, ID: common.WellKnownID_RootCA, DisplayName: "Root CA", ObjectType: NamespaceTypeBuiltInCaRoot},
			{NamespaceID: uuid.Nil, ID: common.WellKnownID_TestRootCA, DisplayName: "Test Intermediate CA", ObjectType: NamespaceTypeBuiltInCaRoot},
		}
	}
*/
func (s *adminServer) ListNamespacesByTypeV2(c *gin.Context, nsType NamespaceTypeShortName) {
	if !authAdminOnly(c) {
		return
	}
	switch nsType {
	case NSTypeRootCA:
		c.JSON(http.StatusOK, getRootCaRefs())
	case NSTypeIntCA:
		c.JSON(http.StatusOK, getIntCaRefs())
	}
	/*
		if namespaceType == NamespaceTypeBuiltInCaRoot {
			c.JSON(http.StatusOK, getBuiltInCaRootNamespaceRefs())
			return
		} else if namespaceType == NamespaceTypeBuiltInCaInt {
			c.JSON(http.StatusOK, getBuiltInCaIntNamespaceRefs())
			return
		}
		if !auth.CallerPrincipalHasAdminRole(c) {
			c.JSON(http.StatusForbidden, gin.H{"message": "only admin can list name spaces"})
			return
		}

		switch namespaceType {
		case NamespaceTypeMsGraphGroup,
			NamespaceTypeMsGraphServicePrincipal,
			NamespaceTypeMsGraphDevice,
			NamespaceTypeMsGraphUser,
			NamespaceTypeMsGraphApplication:
			// allow
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "namespace type not supported"})
			return
		}

		list, err := s.ListDirectoryObjectByType(c, namespaceType)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get list of directory objects")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		results := make([]NamespaceRef, len(list))
		for i, item := range list {
			item.PopulateNamespaceRef(&results[i])
		}
		c.JSON(http.StatusOK, results)
	*/
}

/*
func (s *adminServer) genDirDocFromMsGraph(c *gin.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	dirObj, err := s.msGraphClient.DirectoryObjects().ByDirectoryObjectId(objectID.String()).Get(c, nil)
	if err != nil {
		return nil, err
	}
	doc := new(DirectoryObjectDoc)
	doc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, objectID)
	doc.NamespaceID = wellknownNamespaceID_directoryID
	doc.OdataType = *dirObj.GetOdataType()
	switch doc.OdataType {
	case "#microsoft.graph.user":
		if userObj, ok := dirObj.(msgraphmodels.Userable); ok {
			doc.DisplayName = *userObj.GetDisplayName()
			doc.UserPrincipalName = userObj.GetUserPrincipalName()
		}
	case "#microsoft.graph.servicePrincipal":
		if spObj, ok := dirObj.(msgraphmodels.ServicePrincipalable); ok {
			doc.DisplayName = *spObj.GetDisplayName()
			doc.ServicePrincipalType = spObj.GetServicePrincipalType()
		}
	case "#microsoft.graph.group":
		if gObj, ok := dirObj.(msgraphmodels.Groupable); ok {
			doc.DisplayName = *gObj.GetDisplayName()
		}
	case "#microsoft.graph.device":
		if dObj, ok := dirObj.(msgraphmodels.Deviceable); ok {
			doc.DisplayName = *dObj.GetDisplayName()
			deviceId, err := uuid.Parse(*dObj.GetDeviceId())
			if err != nil {
				return nil, err
			}
			doc.Device = &DirectoryObjectDocDeviceSection{
				DeviceID:               deviceId,
				OperatingSystem:        dObj.GetOperatingSystem(),
				OperatingSystemVersion: dObj.GetOperatingSystemVersion(),
				DeviceOwnership:        dObj.GetDeviceOwnership(),
				IsCompliant:            dObj.GetIsCompliant(),
			}
		}
	case "#microsoft.graph.application":
		if dObj, ok := dirObj.(msgraphmodels.Applicationable); ok {
			doc.DisplayName = *dObj.GetDisplayName()
			doc.AppID = dObj.GetAppId()
		}
	default:
		return nil, fmt.Errorf("graph object type (%s) not supported", doc.OdataType)
	}
	return doc, nil
}

func (s *adminServer) syncDirDoc(c *gin.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	doc, err := s.genDirDocFromMsGraph(c, objectID)
	if err != nil {
		return doc, err
	}

	err = kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, doc)
	if err != nil {
		return doc, err
	}
	uuid.Nil = uuid.New()
	return doc, nil
}

func (s *adminServer) RegisterNamespaceProfile(c *gin.Context, objectID uuid.UUID) (*NamespaceProfile, int, error) {
	doc, err := s.syncDirDoc(c, objectID)

	if err != nil {
		if common.IsGraphODataErrorNotFound(err) {
			return nil, http.StatusNotFound, err
		}
		if common.IsAzNotFound(err) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	nsProfile := new(NamespaceProfile)
	doc.PopulateNamespaceProfile(nsProfile)

	return nsProfile, http.StatusOK, nil
}

func (s *adminServer) RegisterNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin can register namespaces"})
		return
	}
	if namespaceId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No namespace id specified"})
		return
	}
	profile, status, err := s.RegisterNamespaceProfile(c, namespaceId)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.Error().Err(err).Msg("Failed to register graph object")
			c.JSON(status, gin.H{"message": "internal error"})
		} else {
			c.JSON(status, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (s *adminServer) GetNamespaceProfile(c context.Context, namespaceId uuid.UUID) (*NamespaceProfile, error) {
	if IsCANamespace(namespaceId) {
		nsProfile := new(NamespaceProfile)
		nsProfile.NamespaceID = namespaceId
		switch namespaceId {
		case common.WellKnownID_TestRootCA:
			nsProfile.DisplayName = "Test Root CA"
		case common.WellKnownID_RootCA:
			nsProfile.DisplayName = "Root CA"
		case testNamespaceID_IntCA:
			nsProfile.DisplayName = "Test Intermediate CA"
		case wellKnownNamespaceID_IntCaIntranet:
			nsProfile.DisplayName = "Intermediate CA - Intranet"
		case wellKnownNamespaceID_IntCAService:
			nsProfile.DisplayName = "Intermediate CA - Services"
		}
		if IsIntCANamespace(namespaceId) {
			nsProfile.ObjectType = NamespaceTypeBuiltInCaInt
		} else {
			nsProfile.ObjectType = NamespaceTypeBuiltInCaRoot
		}
		return nsProfile, nil
	}
	doc, err := s.getDirectoryObjectDoc(c, namespaceId)
	if common.IsAzNotFound(err) {
		return nil, nil
	}
	nsProfile := new(NamespaceProfile)
	doc.PopulateNamespaceProfile(nsProfile)
	return nsProfile, nil
}

func (s *adminServer) GetNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	// CA profiles are public
	if !IsCANamespace(namespaceId) {
		if _, ok := authNamespaceAdminOrSelf(c, namespaceId); !ok {
			return
		}
	}

	profile, err := s.GetNamespaceProfile(c, namespaceId)
	if err != nil {
		log.Error().Err(err).Msg("Internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, profile)
}
*/