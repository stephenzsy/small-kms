package agentconfig

import (
	"errors"
	"net/http"

	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func ApiProxyGetDockerInfo(c ctx.RequestContext) error {
	return errors.New("Not implemented")
	// nsID := ns.GetNamespaceContext(c).GetID()
	// if !nsID.Identifier().IsUUID() {
	// 	return fmt.Errorf("%w:namespace id is not uuid: %v", common.ErrStatusBadRequest, nsID)
	// }
	// clientPool, ok := c.Value(proxyHttpClientPoolContextKey).(AgentProxyHttpClientPool)
	// if !ok {
	// 	return fmt.Errorf("proxy client pool not found")
	// }

	// client, err := clientPool.GetProxyHttpClient(c.Elevate(), nsID)
	// if err != nil {
	// 	return err
	// }
	// result, err := client.GetDockerInfoWithResponse(c, nsID.Identifier())
	// if err != nil {
	// 	return err
	// }
	// return c.JSON(http.StatusOK, result.JSON200)
}

func ApiGetAgentProxyInfo(c ctx.RequestContext) error {
	nsID := ns.GetNamespaceContext(c).GetID()
	doc := AgentActiveServerCallbackDoc{}
	err := kmsdoc.Read(c, NewAgentCallbackDocLocator(nsID, shared.AgentConfigNameActiveServer), &doc)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, doc.toModel())
}
