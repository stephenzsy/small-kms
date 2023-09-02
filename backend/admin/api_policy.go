package admin

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// returns result with nil id if not found
func (s *adminServer) readPolicyDBItem(c context.Context, namespaceID uuid.UUID, policyID uuid.UUID) (*PolicyDBItem, error) {
	db := s.azCosmosContainerClientPolicies
	resp, err := db.ReadItem(c, azcosmos.NewPartitionKeyString(namespaceID.String()), policyID.String(), nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			// Handle Error
			if respErr.StatusCode == http.StatusNotFound {
				return nil, nil
			}
		}
		return nil, err
	}
	result := new(PolicyDBItem)
	err = json.Unmarshal(resp.Value, result)
	return result, err
}

// returns result with nil id if not found
func (s *adminServer) persistPolicyDBItem(c context.Context, policy *PolicyDBItem) (*PolicyDBItem, error) {
	db := s.azCosmosContainerClientPolicies
	content, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	_, err = db.UpsertItem(c, azcosmos.NewPartitionKeyString(policy.NamespaceID.String()), content, nil)
	return policy, err
}

func (s *adminServer) PutPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	// validate
	callerID, ok := authNamespaceAdminOrSelf(c, namespaceID)
	if !ok {
		return
	}

	p := PolicyDBItem{}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	p.NamespaceID = namespaceID
	p.ID = policyID
	p.UpdatedBy = callerID.String()
	p.UpdatedAt = time.Now().UTC()

	// currently only allow root cert request policy
	if p.Type != PolicyTypeCertRequest || !IsRootCANamespace(namespaceID) {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	policy, err := p.ToCertRequestPolicy()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// write to DB
	policy, err = s.persistPolicyDBItem(c, policy.DBItem())
	if err != nil {
		log.Printf("Internal error: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, policy)
}

func (s *adminServer) GetPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	// validate
	_, ok := authNamespaceAdminOrSelf(c, namespaceID)
	if !ok {
		return
	}
	result, err := s.readPolicyDBItem(c, namespaceID, policyID)
	if err != nil {
		log.Printf("Internal error: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	if result == nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, result)
}
