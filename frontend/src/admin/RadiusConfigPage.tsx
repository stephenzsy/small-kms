import { useContext, useEffect, useMemo } from "react";

import { useRequest } from "ahooks";
import {
  Button,
  Card,
  Checkbox,
  Collapse,
  CollapseProps,
  Form,
  Input,
  Table,
  Typography,
} from "antd";
import { useForm } from "antd/es/form/Form";
import { ColumnsType } from "antd/es/table";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { Link } from "../components/Link";
import {
  AdminApi,
  AgentConfigRadiusFields,
  AgentConfigRadiusToJSON,
  AgentContainerConfiguration,
  AgentInstance,
  NamespaceKind,
  RadiusClientConfig,
  RadiusEapTls,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import {
  NamespaceContext,
  NamespaceContextValue,
} from "./contexts/NamespaceContext";
import {
  RadiusConfigPatchContext,
  useRadiusConfigPatch,
} from "./contexts/RadiusConfigPatchContext";
import { RadiusConfigContainerForm } from "./forms/RadiusConfigContainerForm";
import { RadiusConfigServersForm } from "./forms/RadiusConfigServersForm";

export type AgentServerConfigFormState = AgentContainerConfiguration;

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

type RadiusMiscFormState = Pick<AgentConfigRadiusFields, "debugMode">;
function RadiusMiscForm() {
  const { run, data } = useRadiusConfigPatch();
  const [form] = useForm<RadiusMiscFormState>();
  useEffect(() => {
    if (data && data.debugMode != undefined) {
      form.setFieldsValue({
        debugMode: data.debugMode,
      });
    }
  }, [data, form]);
  return (
    <Form form={form} layout="vertical" onFinish={(v) => run(v)}>
      <Form.Item<RadiusMiscFormState>
        valuePropName="checked"
        getValueFromEvent={(e) => e.target.checked}
        name="debugMode"
      >
        <Checkbox>Debug mode</Checkbox>
      </Form.Item>
      <Form.Item className="mt-6">
        <Button htmlType="submit" type="primary">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

type RadiusEapTlsFormState = Pick<RadiusEapTls, "certPolicyId">;
function RadiusEapTlsForm() {
  const { run, data } = useRadiusConfigPatch();
  const [form] = useForm<RadiusEapTlsFormState>();
  useEffect(() => {
    if (data?.eapTls) {
      form.setFieldsValue(data.eapTls);
    }
  }, [data, form]);
  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(v) =>
        run({
          eapTls: v,
        })
      }
    >
      <Form.Item<RadiusEapTlsFormState>
        name="certPolicyId"
        label="Server TLS certificate policy ID"
        required
      >
        <Input />
      </Form.Item>
      <Form.Item className="mt-6">
        <Button htmlType="submit" type="primary">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

function RadiusClientsForm() {
  const { run, data } = useRadiusConfigPatch();
  const [form] = useForm<Pick<AgentConfigRadiusFields, "clients">>();

  useEffect(() => {
    if (data && data.clients) {
      form.setFieldsValue({
        clients: data.clients,
      });
    }
  }, [data, form]);
  return (
    <Form form={form} layout="vertical" onFinish={(v) => run(v)}>
      <Form.List name={"clients"}>
        {(subFields, subOpt) => {
          return (
            <div className="flex flex-col gap-4">
              {subFields.map((subField) => (
                <div
                  key={subField.key}
                  className=" ring-1 ring-neutral-400 p-4 rounded-md"
                >
                  <Form.Item<RadiusClientConfig[]>
                    name={[subField.name, "name"]}
                    label="Name"
                    required
                  >
                    <Input placeholder={"localhost"} />
                  </Form.Item>
                  <Form.Item<RadiusClientConfig[]>
                    name={[subField.name, "ipaddr"]}
                    label="IP Address or CIDR"
                    required
                  >
                    <Input placeholder={"192.168.0.1/24"} />
                  </Form.Item>
                  <Form.Item<RadiusClientConfig[]>
                    name={[subField.name, "secretPolicyId"]}
                    className="flex-auto"
                    label="Secret Policy ID"
                    required
                  >
                    <Input />
                  </Form.Item>
                  <Button
                    danger
                    onClick={() => {
                      subOpt.remove(subField.name);
                    }}
                  >
                    Remove
                  </Button>
                </div>
              ))}
              <Button type="dashed" onClick={() => subOpt.add()} block>
                Add client configuration
              </Button>
            </div>
          );
        }}
      </Form.List>
      <Form.Item className="mt-6">
        <Button htmlType="submit" type="primary">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

const collapseItems: CollapseProps["items"] = [
  {
    key: "container",
    label: "Container",
    children: <RadiusConfigContainerForm />,
  },
  {
    key: "clients",
    label: "Clients",
    children: <RadiusClientsForm />,
  },
  {
    key: "eap-tls",
    label: "EAP-TLS",
    children: <RadiusEapTlsForm />,
  },
  {
    key: "servers",
    label: "Servers",
    children: <RadiusConfigServersForm />,
  },
  {
    key: "misceallaneous",
    label: "Miscellaneous",
    children: <RadiusMiscForm />,
  },
];

export function AgentInstancesList() {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { data } = useRequest(
    () => {
      return api.listAgentInstances({
        namespaceId: namespaceId,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceKind],
      ready: !!namespaceId && !!namespaceKind,
    }
  );
  const columns = useAgentInstanceColumns(namespaceKind, namespaceId);
  return (
    <Table
      dataSource={data}
      columns={columns}
      rowKey={(r: AgentInstance) => r.id}
    />
  );
}

export default function RadiusConfigPage({
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
  const api = useAuthedClient(AdminApi);
  const { namespaceId, namespaceKind } = nsCtxValue;
  const patchSvc = useRequest(
    (params?: AgentConfigRadiusFields) => {
      if (params) {
        return api.patchAgentConfigRadius({
          agentConfigRadiusFields: params,
          namespaceKind,
          namespaceId,
        });
      }
      return api.getAgentConfigRadius({
        namespaceId,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceKind],
      ready: !!namespaceId && !!namespaceKind,
    }
  );
  const { data: radiusConfig } = patchSvc;

  return (
    <>
      <Typography.Title>
        Radius configuration:{" "}
        {isGlobalConfig ? "Global configuration" : managedApp?.displayName}
      </Typography.Title>

      <NamespaceContext.Provider value={nsCtxValue}>
        <RadiusConfigPatchContext.Provider value={patchSvc}>
          <Card title="Current configuration">
            <JsonDataDisplay
              data={radiusConfig}
              toJson={AgentConfigRadiusToJSON}
            />
          </Card>
          <Card title="RADIUS configuration">
            <Collapse items={collapseItems} />
          </Card>
          {!isGlobalConfig && (
            <Card title="Instances">
              <AgentInstancesList />
            </Card>
          )}
        </RadiusConfigPatchContext.Provider>
      </NamespaceContext.Provider>
    </>
  );
}
