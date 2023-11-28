import { useRequest } from "ahooks";
import { Result } from "ahooks/lib/useRequest/src/types";
import {
  Button,
  Card,
  Form,
  TableColumnType,
  Tabs,
  TabsProps,
  Typography,
} from "antd";
import { useContext, useEffect, useMemo } from "react";
import { ResourceRefsSelect } from "../../admin/ResourceRefsSelect";
import { AgentContext } from "../../admin/contexts/AgentContext";
import { NamespaceContext } from "../../admin/contexts/NamespaceContext";
import { useNamespace } from "../../admin/contexts/useNamespace";
import { ResourceRefsTable } from "../../admin/tables/ResourceRefsTable";
import { Link } from "../../components/Link";
import {
  AdminApi,
  AgentConfigEndpoint,
  AgentConfigEndpointFields,
  AgentConfigIdentity,
  AgentConfigIdentityFields,
  AgentConfigName,
  AgentInstanceRef,
  CreateAgentConfigRequest,
  NamespaceProvider,
  Ref as ResourceRef,
} from "../../generated/apiv2";
import { useAdminApi } from "../../utils/useCertsApi";

export default function AgentPage() {
  const { agent } = useContext(AgentContext);

  const tabItems = useMemo(
    (): TabsProps["items"] => [
      {
        key: "dashboard",
        label: "Dashboard",
        children: <AgentDashboardCards />,
      },
      {
        key: "configurations",
        label: "Configurations",
        children: <AgentConfigurationsCards />,
      },
    ],
    []
  );

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
      {agent?.servicePrincipalId && (
        <NamespaceContext.Provider
          value={{
            namespaceId: agent.servicePrincipalId,
            namespaceKind: NamespaceProvider.NamespaceProviderServicePrincipal,
          }}
        >
          <Tabs destroyInactiveTabPane items={tabItems} />
        </NamespaceContext.Provider>
      )}
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

function AgentDashboardCards() {
  const { namespaceId } = useNamespace();
  const api = useAdminApi();
  const { data: instances } = useRequest(
    async () => {
      const instances = await api?.listAgentInstances({
        namespaceId: namespaceId,
      });

      // const i = instances?.[0];
      // api?.getAgentDiagnostics({
      //   namespaceId: namespaceId,
      //   id: i?.id as string,
      // });

      return instances;
    },
    { refreshDeps: [namespaceId] }
  );

  const extraColumns = useMemo(
    (): TableColumnType<AgentInstanceRef>[] => [
      {
        title: "Endpoint",
        dataIndex: "endpoint",
        key: "endpoint",
      },
      {
        title: "State",
        dataIndex: "state",
        key: "state",
        render: (state: AgentInstanceRef["state"]) => {
          return <span className="capitalize">{state}</span>;
        },
      },
      {
        title: "Actions",
        key: "actions",
        render: (instance: AgentInstanceRef) => {
          return (
            <Link to={`./instances/${instance.id}/dashboard`}>Dashboard</Link>
          );
        },
      },
    ],
    []
  );

  return (
    <div className="space-y-4">
      <Card title="Instances">
        <ResourceRefsTable
          dataSource={instances}
          noDisplayName
          extraColumns={extraColumns}
        />
      </Card>
    </div>
  );
}

function AgentConfigurationsCards() {
  const api = useAdminApi();
  const { namespaceId } = useNamespace();
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
    <div className="space-y-4">
      <Card title="Identity">
        <AgentIdentityForm
          onPutConfig={updateIdentityConfig}
          certificatePolicies={certificatePoliciesResult}
          value={identityConfig}
        />
      </Card>
      <Card title="Endpoint">
        <AgentEndpointForm
          certificatePolicies={certificatePoliciesResult}
          keyPolicies={keyPoliciesResult}
          onPutConfig={updateEndpointConfig}
          value={endpointConfig}
        />
      </Card>
    </div>
  );
}
