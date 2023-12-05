import { useContext, useEffect, useMemo } from "react";

import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Form, Input, Table, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import { ColumnsType } from "antd/es/table";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { Link } from "../components/Link";
import {
  AdminApi,
  AgentConfigServerFields,
  AgentConfigServerToJSON,
  AgentInstance,
  NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import {
  NamespaceContext,
  NamespaceContextValue,
} from "./contexts/NamespaceContext";

// const selectOptions: Array<SelectItem<AgentConfigName>> = [
//   {
//     id: AgentConfigName.AgentConfigNameActiveHostBootstrap,
//     name: "Agent Active Host Bootstrap",
//   },
//   {
//     id: AgentConfigName.AgentConfigNameActiveServer,
//     name: "Agent Active Host Server",
//   },
// ];

type AgentServerConfigFormState = Partial<AgentConfigServerFields>;

function AgentConfigServerFormCard({
  isGlobalConfig,
}: {
  isGlobalConfig: boolean;
}) {
  const [form] = useForm<AgentServerConfigFormState>();
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { data: agentServerConfig, run } = useRequest(
    (params?: Partial<AgentConfigServerFields>) => {
      if (params) {
        return api.putAgentConfigServer({
          agentConfigServerFields: params as AgentConfigServerFields,
          namespaceKind,
          namespaceId: namespaceIdentifier,
        });
      }
      return api.getAgentConfigServer({
        namespaceId: namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier && !!namespaceKind,
    }
  );

  useEffect(() => {
    if (agentServerConfig) {
      form.setFieldsValue(agentServerConfig);
    }
  }, [agentServerConfig, form]);

  const setCurrentBuildTag = useMemoizedFn(async () => {
    // const tag = (await api.getDiagnostics()).serviceRuntime.buildId.split(
    //   "."
    // )[0];
    // const currentValue = form.getFieldValue("azureAcrImageRef");
    // const currentPrfix = currentValue.split(":")[0];
    // form.setFieldValue("azureAcrImageRef", `${currentPrfix}:${tag}`);
  });

  return (
    <Card title="Agent server configuration">
      <div className="mb-6">
        Current configuration:
        <JsonDataDisplay
          data={agentServerConfig}
          toJson={AgentConfigServerToJSON}
        />
      </div>
      <Form
        form={form}
        layout="vertical"
        onFinish={(values) => {
          run(values);
        }}
      >
        {isGlobalConfig && (
          <Form.Item<AgentServerConfigFormState>
            name="azureAcrImageRef"
            label={
              <span>
                Azure Container Registry image Reference{" "}
                <Button type="link" onClick={setCurrentBuildTag}>
                  Use current build tag
                </Button>
              </span>
            }
            required
          >
            <Input placeholder="example.com/image:latest" />
          </Form.Item>
        )}
        {!isGlobalConfig && (
          <Form.Item<AgentServerConfigFormState>
            name="jwtKeyCertPolicyId"
            label="Json web token key certificate policy"
            required
          >
            <Input />
          </Form.Item>
        )}
        <Form.Item>
          <Button htmlType="submit" type="primary">
            Submit
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
}

function useAgentInstanceColumns(
  namespaceKind: NamespaceKind,
  namespaceIdentifier: string
) {
  return useMemo(
    (): ColumnsType<AgentInstance> => [
      {
        key: "id",
        title: "ID",
        render: (r: AgentInstance) => <span className="font-mono">{r.id}</span>,
      },
      {
        key: "endpoint",
        title: "Endpoint",
        render: (r: AgentInstance) => (
          <span className="font-mono">{r.endpoint}</span>
        ),
      },
      {
        key: "version",
        title: "Config version",
        render: (r: AgentInstance) => (
          <span className="font-mono">{r.version}</span>
        ),
      },
      {
        key: "buildID",
        title: "Build ID",
        render: (r: AgentInstance) => (
          <span className="font-mono">{r.buildId}</span>
        ),
      },
      {
        key: "actions",
        title: "Actions",
        render: (r: AgentInstance) => (
          <Link
            to={`/app/${namespaceKind}/${namespaceIdentifier}/agent/${r.id}/dashboard`}
          >
            Dashboard
          </Link>
        ),
      },
    ],
    [namespaceKind, namespaceIdentifier]
  );
}

export function AgentInstancesList() {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { data } = useRequest(
    () => {
      return api.listAgentInstances({
        namespaceId: namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier && !!namespaceKind,
    }
  );
  const columns = useAgentInstanceColumns(namespaceKind, namespaceIdentifier);
  return (
    <Table
      dataSource={data}
      columns={columns}
      rowKey={(r: AgentInstance) => r.id}
    />
  );
}

export default function ProvisionAgentPage({
  isGlobalConfig = false,
}: {
  isGlobalConfig?: boolean;
}) {
  const { managedApp } = useContext(ManagedAppContext);
  const nsCtxValue: NamespaceContextValue = useMemo(
    () =>
      isGlobalConfig
        ? {
            namespaceKind: NamespaceKind.NamespaceKindSystem,
            namespaceId: "default",
          }
        : {
            namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
            namespaceId: managedApp?.servicePrincipalId ?? "",
          },
    [isGlobalConfig, managedApp]
  );
  return (
    <>
      <Typography.Title>
        Provision agent:{" "}
        {isGlobalConfig ? "Global configuration" : managedApp?.displayName}
      </Typography.Title>

      <NamespaceContext.Provider value={nsCtxValue}>
        <AgentConfigServerFormCard isGlobalConfig={isGlobalConfig} />
        {!isGlobalConfig && (
          <Card title="Instances">
            <AgentInstancesList />
          </Card>
        )}
      </NamespaceContext.Provider>
    </>
  );
}
