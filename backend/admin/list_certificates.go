package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *adminServer) listCertItems(ctx context.Context, namespaceID uuid.UUID) (results []CertDBItem, err error) {
	db := s.config.AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(namespaceID.String())
	pager := db.NewQueryItemsPager(`SELECT * FROM c
		WHERE c.namespaceId = @namespaceId
		ORDER BY c.notAfter DESC`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: namespaceID.String()},
			},
		})

	for pager.More() {
		t, scanErr := pager.NextPage(ctx)
		if scanErr != nil {
			err = fmt.Errorf("faild to get list of certificates: %w", scanErr)
			return
		}
		for _, itemBytes := range t.Items {
			item := CertDBItem{}
			if err = json.Unmarshal(itemBytes, &item); err != nil {
				err = fmt.Errorf("faild to serialize db entry: %w", err)
				return
			}
			results = append(results, item)
		}
	}
	return
}

func (s *adminServer) ListCertificatesV1(c *gin.Context, namespaceId uuid.UUID) {
	items, err := s.listCertItems(c, namespaceId)
	if err != nil {
		log.Printf("Faild to get list of certificates: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	results := make([]CertificateRef, len(items))
	for i, item := range items {
		results[i] = item.CertificateRef
	}
	c.JSON(200, results)
}
