package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

type CertDBItem struct {
	CertificateRef
	KeyStore  string `json:"keyStore,omitempty"`
	CertStore string `json:"certStore,omitempty"`
}

func (s *adminServer) findLatestCertificate(ctx context.Context, namespaceID NamespaceID, name string) (result CertDBItem, err error) {
	partitionKey := azcosmos.NewPartitionKeyString(namespaceID.String())
	db := s.azCosmosContainerClientCerts
	pager := db.NewQueryItemsPager(`
SELECT TOP 1
	*
FROM c
WHERE c.namespaceId = @namespaceId AND c.name = @name
ORDER BY c.notAfter DESC`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: namespaceID},
				{Name: "@name", Value: name},
			},
		})
	t, err := pager.NextPage(ctx)
	if err != nil {
		return
	}
	if len(t.Items) > 0 {
		err = json.Unmarshal(t.Items[0], &result)
	}
	return
}

// returns result with nil id if not found
func (s *adminServer) ReadCertDBItem(c context.Context, namespaceID NamespaceID, id uuid.UUID) (result CertDBItem, err error) {
	db := s.azCosmosContainerClientCerts
	resp, err := db.ReadItem(c, azcosmos.NewPartitionKeyString(namespaceID.String()), id.String(), nil)
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
	return
}

type CertificateEnrollmentPolicyDTO struct {
	CertificateEnrollmentPolicy
	ValidityDuration time.Duration `json:"-"`
}
