package agentconfig

import (
	"fmt"
	"net/http"

	"github.com/stephenzsy/small-kms/backend/common"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func ApiProxyGetDockerInfo(c common.RequestContext) error {
	nsID := ns.GetNamespaceContext(c).GetID()
	if !nsID.Identifier().IsUUID() {
		return fmt.Errorf("%w:namespace id is not uuid: %v", common.ErrStatusBadRequest, nsID)
	}
	clientPool, ok := c.Value(proxyHttpClientPoolContextKey).(AgentProxyHttpClientPool)
	if !ok {
		return fmt.Errorf("proxy client pool not found")
	}

	client, err := clientPool.GetProxyHttpClient(c.Elevate(), nsID)
	if err != nil {
		return err
	}
	result, err := client.GetDockerInfoWithResponse(c, nsID.Identifier())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result.JSON200)
}
