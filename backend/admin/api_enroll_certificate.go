package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
)

// extract owner name

const (
	deviceOwnershipTypeCompany  = "Company"
	deviceOwnershipTypePersonal = "Personal"
)

func (s *adminServer) validateEnrollCertificateRequest(c *gin.Context, r *CertificateEnrollRequest) (int, error) {
	// validate device namespace
	profile, status, err := s.RegisterNamespaceProfile(c, r.OwnerNamespaceID)
	if err != nil {
		return status, err
	}
	if profile.ObjectType != NamespaceTypeMsGraphDevice {
		return http.StatusBadRequest, fmt.Errorf("owner namespace is not for device: %s", r.OwnerNamespaceID)
	}
	if profile.DeviceOwnership == nil {
		return http.StatusBadRequest, fmt.Errorf("device must be registered as Company or Personal: %s", profile.ID)
	}
	if profile.IsCompliant == nil || !*profile.IsCompliant {
		return http.StatusBadRequest, fmt.Errorf("device must be compliant: %s", profile.ID)
	}

	verifyIds := make([]string, 1, 2)
	verifyIds[0] = profile.ID.String()

	switch *profile.DeviceOwnership {
	case deviceOwnershipTypeCompany:
		// TODO company owned, admin can manage or owner can manage
	case deviceOwnershipTypePersonal:
		// verify owner
		if r.DeviceOwnerID == nil {
			return http.StatusBadRequest, fmt.Errorf("null device owner id: %s", profile.ID)
		}

		// check graph
		roResp, err := s.msGraphClient.Devices().ByDeviceId(*profile.DeviceID).RegisteredOwners().Get(c, &devices.ItemRegisteredOwnersRequestBuilderGetRequestConfiguration{
			QueryParameters: &devices.ItemRegisteredOwnersRequestBuilderGetQueryParameters{
				Select: []string{"id"},
			},
		})
		if err != nil {
			return http.StatusInternalServerError, err
		}
		isOwner := false
		for _, v := range roResp.GetValue() {
			if parsedOwnerId, err := uuid.Parse(*v.GetId()); err == nil && parsedOwnerId == *r.DeviceOwnerID {
				isOwner = true
				break
			}
		}
		if !isOwner {
			return http.StatusBadRequest, fmt.Errorf("device owner id does not match: %s, owner: %s", *profile.DeviceID, r.DeviceOwnerID)
		}
		// verify ids belong to group
		verifyIds = append(verifyIds, r.DeviceOwnerID.String())
	default:
		return http.StatusBadRequest, fmt.Errorf("unsupported device ownership: %s [%s]", *profile.DeviceOwnership, profile.ID)
	}

	// verify group membership
	groupProfile, status, err := s.RegisterNamespaceProfile(c, *r.PolicyNamespaceID)
	if err != nil {
		return status, err
	}
	if groupProfile.ObjectType != NamespaceTypeMsGraphGroup {
		return http.StatusBadRequest, fmt.Errorf("policy namespace is not for group: %s", r.PolicyNamespaceID)
	}

	checkMembershipRequestBody := directoryobjects.NewItemCheckMemberObjectsPostRequestBody()
	checkMembershipRequestBody.SetIds(verifyIds)

	checkMemberResp, err := s.msGraphClient.DirectoryObjects().ByDirectoryObjectId(groupProfile.ID.String()).CheckMemberObjects().Post(c, checkMembershipRequestBody, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(verifyIds) != len(checkMemberResp.GetValue()) {
		return http.StatusBadRequest, fmt.Errorf("not all ids are members of group: %s, memberIDs: %s", groupProfile.ID, strings.Join(verifyIds, ","))
	}
	return 0, nil
}

func (s *adminServer) EnrollCertificateV1(c *gin.Context) {
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin can enroll certificate"})
		return
	}
	req := CertificateEnrollRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validate request
	if statusCode, err := s.validateEnrollCertificateRequest(c, &req); err != nil {
		if statusCode >= 500 {
			log.Error().Err(err).Msg("Failed to validate enroll certificate request")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		} else {
			c.JSON(statusCode, gin.H{"message": err.Error()})
		}
		return
	}
}
