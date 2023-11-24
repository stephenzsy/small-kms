import { useRequest } from "ahooks";
import { Result } from "ahooks/lib/useRequest/src/types";
import { Button, Card, Divider, Form, Typography } from "antd";
import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { ResourceRefsSelect } from "../../admin/ResourceRefsSelect";
import { NamespaceContext } from "../../admin/contexts/NamespaceContext";
import { Link } from "../../components/Link";
import {
  AdminApi,
  AgentConfigEndpoint,
  AgentConfigEndpointFields,
  AgentConfigIdentity,
  AgentConfigIdentityFields,
  AgentConfigName,
  CreateAgentConfigRequest,
  NamespaceProvider,
  Ref as ResourceRef,
} from "../../generated/apiv2";
import { useAdminApi, useAuthedClientV2 } from "../../utils/useCertsApi";

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
            <dt>ID (Client ID)</dt>
            <dd className="font-mono">{agent?.id}</dd>
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
            <dd className="dd font-mono">
              <Link to={`/service-principal/${agent?.servicePrincipalId}`}>
                {agent?.servicePrincipalId}
              </Link>
            </dd>
          </div>
        </dl>
      </Card>
      <NamespaceContext.Provider
        value={{
          namespaceId: agent?.servicePrincipalId ?? "",
          namespaceKind: NamespaceProvider.NamespaceProviderServicePrincipal,
        }}
      >
        {agent?.servicePrincipalId && (
          <AgentConfigurationsCard namespaceId={agent.servicePrincipalId} />
        )}
      </NamespaceContext.Provider>
    </>
  );
}

function AgentIdentityForm({
  certificatePolicies,
  value,
  onPutConfig,
}: {
  value: AgentConfigIdentity | undefined;
  certificatePolicies: Result<ResourceRef[] | undefined, []>;

  onPutConfig: (config: AgentConfigIdentityFields) => void;
}) {
  const [formInstance] = Form.useForm<AgentConfigIdentityFields>();

  useEffect(() => {
    if (value) {
      formInstance.setFieldsValue(value);
    }
  }, [value, formInstance]);

  return (
    <Form form={formInstance} layout="vertical" onFinish={onPutConfig}>
      <Form.Item<AgentConfigIdentityFields>
        name="keyCredentialCertificatePolicyId"
        label="Select Key Credential Certificate Policy"
      >
        <ResourceRefsSelect
          data={certificatePolicies.data}
          loading={certificatePolicies.loading}
        />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

function AgentEndpointForm({
  certificatePolicies,
  keyPolicies,
  value,
  onPutConfig,
}: {
  value: AgentConfigEndpoint | undefined;
  certificatePolicies: Result<ResourceRef[] | undefined, []>;
  keyPolicies: Result<ResourceRef[] | undefined, []>;
  onPutConfig: (config: AgentConfigEndpointFields) => void;
}) {
  const [formInstance] = Form.useForm<AgentConfigEndpointFields>();

  useEffect(() => {
    if (value) {
      formInstance.setFieldsValue(value);
    }
  }, [value, formInstance]);

  return (
    <Form form={formInstance} layout="vertical" onFinish={onPutConfig}>
      <Form.Item<AgentConfigEndpointFields>
        name="tlsCertificatePolicyId"
        label="Select TLS Certificate Policy"
      >
        <ResourceRefsSelect
          data={certificatePolicies.data}
          loading={certificatePolicies.loading}
        />
      </Form.Item>
      <Form.Item<AgentConfigEndpointFields>
        name="jwtVerifyKeyPolicyId"
        label="Select JWT Verification Key Policy"
      >
        <ResourceRefsSelect
          data={keyPolicies.data}
          loading={keyPolicies.loading}
        />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

function useAgentConfigRequest<
  T extends AgentConfigIdentity | AgentConfigEndpoint
>(api: AdminApi | undefined, namespaceId: string, name: AgentConfigName) {
  return useRequest(
    async (req?: CreateAgentConfigRequest): Promise<T> => {
      if (req) {
        return (await api?.putAgentConfig({
          configName: name,
          namespaceId: namespaceId,
          createAgentConfigRequest: req,
        })) as T;
      }
      return (await api?.getAgentConfig({
        configName: name,
        namespaceId: namespaceId,
      })) as T;
    },
    {
      refreshDeps: [namespaceId, name],
    }
  );
}

function AgentConfigurationsCard({ namespaceId }: { namespaceId: string }) {
  const api = useAdminApi();
  const { run: updateIdentityConfig, data: identityConfig } =
    useAgentConfigRequest<AgentConfigIdentity>(
      api,
      namespaceId,
      AgentConfigName.AgentConfigNameIdentity
    );
  const { run: updateEndpointConfig, data: endpointConfig } =
    useAgentConfigRequest<AgentConfigEndpoint>(
      api,
      namespaceId,
      AgentConfigName.AgentConfigNameEndpoint
    );
  const certificatePoliciesResult = useRequest(
    async () => {
      return api?.listCertificatePolicies({
        namespaceId,
        namespaceProvider: NamespaceProvider.NamespaceProviderServicePrincipal,
      });
    },
    { refreshDeps: [namespaceId] }
  );
  const keyPoliciesResult = useRequest(
    async () => {
      return api?.listKeyPolicies({
        namespaceId,
        namespaceProvider: NamespaceProvider.NamespaceProviderServicePrincipal,
      });
    },
    { refreshDeps: [namespaceId] }
  );
  return (
    <Card title="Agent configurations">
      <section>
        <Typography.Title level={4}>Identity</Typography.Title>
        <AgentIdentityForm
          onPutConfig={updateIdentityConfig}
          certificatePolicies={certificatePoliciesResult}
          value={identityConfig}
        />
      </section>
      <Divider />
      <section>
        <Typography.Title level={4}>Endpoint</Typography.Title>
        <AgentEndpointForm
          certificatePolicies={certificatePoliciesResult}
          keyPolicies={keyPoliciesResult}
          onPutConfig={updateEndpointConfig}
          value={endpointConfig}
        />
      </section>
    </Card>
  );
}
