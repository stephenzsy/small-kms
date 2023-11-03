import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Typography } from "antd";
import { useParams } from "react-router-dom";
import { AdminApi, PullImageRequest } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { useContext } from "react";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { useForm } from "antd/es/form/Form";

type DockerPullImageFormState = {
  imageTag: string;
};

function DockerPullImageForm({
  instanceId,
  token,
}: {
  instanceId: string;
  token: string;
}) {
  const [form] = useForm<DockerPullImageFormState>();
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const api = useAuthedClient(AdminApi);

  const { run: pullImage } = useRequest(
    async (req: PullImageRequest) => {
      await api.agentPullImage({
        namespaceIdentifier,
        namespaceKind,
        resourceIdentifier: instanceId,
        pullImageRequest: req,
        xCryptocatProxyAuthorization: token,
      });
    },
    { manual: true }
  );

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(s) => {
        pullImage({
          imageTag: s.imageTag,
        });
      }}
    >
      <Form.Item<DockerPullImageFormState> name="imageTag" label="Image tag">
        <Input placeholder="latest" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

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

  const { data: dockerImages, run: getDockerImages } = useRequest(
    async () => {
      if (instanceId && tokenResult?.accessToken) {
        return await api.agentDockerImageList({
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
          xCryptocatProxyAuthorization: tokenResult?.accessToken,
        });
      }
    },
    { manual: true }
  );

  const { data: dockerContainers, run: listDockerContainers } = useRequest(
    async () => {
      if (instanceId && tokenResult?.accessToken) {
        return await api.agentDockerContainerList({
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
          xCryptocatProxyAuthorization: tokenResult?.accessToken,
        });
      }
    },
    { manual: true }
  );

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

      <Card title="Docker images">
        <Button
          type="primary"
          onClick={getDockerImages}
          disabled={!tokenResult}
        >
          List Docker images
        </Button>
        <JsonDataDisplay data={dockerImages} />
      </Card>
      {instanceId && tokenResult && (
        <Card title="Docker pull image">
          <DockerPullImageForm
            instanceId={instanceId}
            token={tokenResult.accessToken}
          />
        </Card>
      )}

      <Card title="Docker conatiners">
        <Button
          type="primary"
          onClick={listDockerContainers}
          disabled={!tokenResult}
        >
          List Docker containers
        </Button>
        <JsonDataDisplay data={dockerContainers} />
      </Card>
    </>
  );
}
