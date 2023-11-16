import { useRequest } from "ahooks";
import { AdminApi } from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useParams } from "react-router-dom";
import { Card, Divider, Typography } from "antd";

function useAgent(agentId: string | undefined) {
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    async () => {
      if (agentId) {
        return api.getAgent({
          id: agentId,
        });
      }
    },
    {
      refreshDeps: [agentId],
    }
  );
}

export default function AgentPage() {
  const { id } = useParams<{ id: string }>();
  const { data: agent } = useAgent(id);
  return (
    <>
      <Typography.Title>Agent</Typography.Title>
      <Card title="Agent info">
        <dl className="dl">
          <div className="dl-item">
            <dt className="dt">ID (Client ID)</dt>
            <dd className="dd font-mono">{agent?.id}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Display name</dt>
            <dd className="dd">{agent?.displayName}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Application ID (Object ID)</dt>
            <dd className="dd font-mono">{agent?.applicationId}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Service Principal ID (Object ID)</dt>
            <dd className="dd font-mono">{agent?.servicePrincipalId}</dd>
          </div>
        </dl>
      </Card>
      <Divider />
      <Typography.Title level={2}>Service Principal</Typography.Title>
    </>
  );
}
