package admin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) ListNamespacesV1(c *gin.Context, namespaceType NamespaceType) {
	/*
		if namespaceType == NamespaceTypeBuiltInCaRoot {
			c.JSON(http.StatusOK, nil)
			return
		} else if namespaceType == NamespaceTypeBuiltInCaInt {
			c.JSON(http.StatusOK, nil)
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

		list, err := s.listDirectoryObjectByType(c, string(namespaceType))
		if err != nil {
			log.Error().Err(err).Msg("Failed to get list of directory objects")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		results := make([]NamespaceRef, len(list))
		for i, item := range list {
			item.PopulateNamespaceRef(&results[i])
		}
	*/
	c.JSON(http.StatusOK, nil)
}

func (s *adminServer) genDirDocFromMsGraph(c context.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	return nil, nil
}

// Deprecated: this operation can throw error with graph 404
func (s *adminServer) syncDirDoc(c context.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	doc, err := s.genDirDocFromMsGraph(c, objectID)
	if err != nil {
		return doc, err
	}

	err = kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc)
	if err != nil {
		return doc, err
	}
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

func (s *adminServer) GetNamespaceProfile(c context.Context, namespaceId uuid.UUID) (*NamespaceProfile, error) {
	if isAllowedCaNamespace(namespaceId) {
		nsProfile := new(NamespaceProfile)
		nsProfile.NamespaceID = namespaceId
		switch namespaceId {
		case common.WellKnownID_TestRootCA:
			nsProfile.DisplayName = "Test Root CA"
		case common.WellKnownID_RootCA:
			nsProfile.DisplayName = "Root CA"
		case common.WellKnownID_TestIntCA:
			nsProfile.DisplayName = "Test Intermediate CA"
		case common.WellKnownID_IntCAIntranet:
			nsProfile.DisplayName = "Intermediate CA - Intranet"
		case common.WellKnownID_IntCAService:
			nsProfile.DisplayName = "Intermediate CA - Services"
		}
		if isAllowedIntCaNamespace(namespaceId) {
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
	if !isAllowedCaNamespace(namespaceId) {
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
