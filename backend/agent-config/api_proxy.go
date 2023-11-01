package agentconfig

import (
	"errors"

	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func ApiProxyGetDockerInfo(c ctx.RequestContext) error {
	return errors.New("not implemented")
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
