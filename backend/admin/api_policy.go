package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) ListPoliciesV1(c *gin.Context, namespaceID uuid.UUID) {
	// validate
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}
	results := make([]PolicyRef, 0)
	c.JSON(http.StatusOK, results)
}

var (
	defaultPolicyIdCertEnroll = common.GetID(common.DefaultPolicyIdCertEnroll)
)

func resolvePolicyIdentifier(policyIdentifier string) (uuid.UUID, error) {
	switch policyIdentifier {
	case string(PolicyTypeCertEnroll):
		return defaultPolicyIdCertEnroll, nil
	}
	return uuid.Parse(policyIdentifier)
}

func (s *adminServer) GetPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyIdentifier string) {
	// validate
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}

	c.JSON(200, nil)
}
