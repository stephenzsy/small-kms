import { useRequest } from "ahooks";
import { Card, Typography } from "antd";
import { useParams } from "react-router-dom";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

export default function AgentDashboardPage() {
  const { namespaceId } = useParams() as { namespaceId: string };

  const client = useAuthedClient(AdminApi);
  const { run } = useRequest(
    async () => {
      await client.getDockerInfo({
        namespaceId,
      });
    },
    { manual: true }
  );

  const { data: proxyInfo } = useRequest(() => {
    return client.getAgentProxyInfo({
      namespaceId,
    });
  }, {});
  return (
    <>
      <Typography.Title>Agent Dashboard</Typography.Title>
      <Card title="Agent proxy information">
        <pre>{JSON.stringify(proxyInfo, undefined, 2)}</pre>
      </Card>
    </>
  );
}
