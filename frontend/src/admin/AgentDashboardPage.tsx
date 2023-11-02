import { useRequest } from "ahooks";
import { Button, Card, Typography } from "antd";
import { useParams } from "react-router-dom";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { useContext } from "react";
import { JsonDataDisplay } from "../components/JsonDataDisplay";

export default function AgentDashboardPage() {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const { instanceId } = useParams<{ instanceId: string }>();

  const api = useAuthedClient(AdminApi);
  const { data } = useRequest(
    async () => {
      if (instanceId) {
        return await api.getAgentInstance({
          namespaceKind,
          namespaceIdentifier,
          resourceIdentifier: instanceId,
        });
      }
    },
    {
      refreshDeps: [namespaceKind, namespaceIdentifier, instanceId],
    }
  );

  const { data: tokenResult, run: acquireToken } = useRequest(
    async () => {
      if (instanceId) {
        return await api.createAgentInstanceProxyAuthToken({
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
        });
      }
    },
    {
      manual: true,
    }
  );

  const { data: agentDiag, run: getAgentDiagnostics } = useRequest(
    async () => {
      if (instanceId && tokenResult?.accessToken) {
        return await api.getAgentDiagnostics({
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
          xCryptocatProxyAuthorization: tokenResult?.accessToken,
        });
      }
    },
    { manual: true }
  );

  const { data: dockerInfo, run: getDockerInfo } = useRequest(
    async () => {
      if (instanceId && tokenResult?.accessToken) {
        return await api.getAgentDockerInfo({
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
          xCryptocatProxyAuthorization: tokenResult?.accessToken,
        });
      }
    },
    { manual: true }
  );

  console.log(agentDiag)
  return (
    <>
      <Typography.Title>Agent Dashboard</Typography.Title>
      <Card title="Agent proxy information">
        <JsonDataDisplay data={data} />
        <div className="mt-6">
          <Button type="primary" onClick={acquireToken}>
            Authorize
          </Button>
        </div>
      </Card>
      <Card title="Diagnostics">
        <Button
          type="primary"
          onClick={getAgentDiagnostics}
          disabled={!tokenResult}
        >
          Get Diagnostics
        </Button>
        <JsonDataDisplay data={agentDiag} />
      </Card>
      <Card title="Docker info">
        <Button type="primary" onClick={getDockerInfo} disabled={!tokenResult}>
          Get Docker info
        </Button>
        <JsonDataDisplay data={dockerInfo} />
      </Card>
    </>
  );
}
