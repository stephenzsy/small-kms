package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
)

func (s *adminServer) putPolicyCertEnroll(c *gin.Context, namespaceID NamespaceID, p CertificateEnrollmentParameters) (dto CertificateEnrollmentPolicyDTO, status int, err error) {
	status = 500

	// validate validity period
	if dto.ValidityDuration, err = time.ParseDuration(p.Validity); err != nil {
		status = 400
		return
	}
	dto.CertificateEnrollmentPolicy = CertificateEnrollmentPolicy{
		Validity:         dto.ValidityDuration.String(),
		DelegatedService: p.DelegatedService,
		PolicyID:         PolicyIDCertEnroll,
		IssuerID:         p.IssuerID,
		KeyParameters:    p.KeyParameters,
		NamespaceID:      namespaceID,
		UpdatedAt:        time.Now().UTC(),
		UpdatedBy:        auth.GetCallerID(c),
	}

	// check issuer
	issuer, err := s.ReadCertDBItem(c, dto.NamespaceID, dto.IssuerID)
	if err != nil {
		err = fmt.Errorf("failed to read issuer: %w", err)
		return
	}
	if issuer.ID == uuid.Nil {
		err = fmt.Errorf("issuer not found: %s/%s", dto.NamespaceID.String(), dto.IssuerID.String())
		status = 400
		return
	}
	if issuer.NotAfter.Before(time.Now().UTC().Add(time.Hour * 24 * 30)) {
		err = fmt.Errorf("issuer expiring soon: %s/%s", dto.NamespaceID.String(), dto.IssuerID.String())
		status = 400
		return
	}

	// write to db
	db := s.azCosmosContainerClientPolicies
	itemContent, err := json.Marshal(dto)
	if err != nil {
		return
	}
	_, err = db.UpsertItem(c, azcosmos.NewPartitionKeyString(namespaceID.String()), itemContent, nil)
	if err != nil {
		err = fmt.Errorf("failed to persist policy: %w", err)
		return
	}

	status = 200
	return
}

// returns result with nil id if not found
func (s *adminServer) ReadCertEnrollPolicyDBItem(c context.Context, namespaceID NamespaceID) (result CertificateEnrollmentPolicyDTO, err error) {
	db := s.azCosmosContainerClientPolicies
	resp, err := db.ReadItem(c, azcosmos.NewPartitionKeyString(namespaceID.String()), string(PolicyIDCertEnroll), nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			// Handle Error
			if respErr.StatusCode == http.StatusNotFound {
				return result, nil
			}
		}
		return
	}
	err = json.Unmarshal(resp.Value, &result)
	if err != nil {
		return
	}
	if result.ValidityDuration, err = time.ParseDuration(result.Validity); err != nil {
		return
	}
	return
}

func (s *adminServer) PutPolicyCertEnrollV1(c *gin.Context, namespaceID NamespaceID) {
	// validate
	if !auth.CallerHasAdminAppRole(c) {
		c.JSON(403, gin.H{"error": "App.Admin role required"})
		return
	}

	params := CertificateEnrollmentParameters{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var result CertificateEnrollmentPolicyDTO
	status := 500
	var err error
	switch namespaceID {
	case wellKnownNamespaceID_IntCaSCEPIntranet:
		result, status, err = s.putPolicyCertEnroll(c, namespaceID, params)
	default:
		c.JSON(400, gin.H{"error": fmt.Sprintf("policy not allowed in the specified namespace : %s", namespaceID.String())})
		return
	}
	if status >= 500 {
		log.Printf("Internal error: %s", err.Error())
		c.JSON(status, gin.H{"error": "internal error"})
		return
	}
	if status >= 400 {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (s *adminServer) GetPolicyCertEnrollV1(c *gin.Context, namespaceID NamespaceID) {
	// validate

	if !auth.CallerHasAdminAppRole(c) {
		c.JSON(403, gin.H{"error": "App.Admin role required"})
		return
	}
	result, err := s.ReadCertEnrollPolicyDBItem(c, namespaceID)
	if err != nil {
		log.Printf("Internal error: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	if len(result.PolicyID) == 0 {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, result)
}
