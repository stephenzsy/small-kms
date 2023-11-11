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
import { DefaultOptionType } from "antd/es/select";
import { ColumnsType } from "antd/es/table";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { Link } from "../components/Link";
import {
  AdminApi,
  AgentConfigName,
  AgentConfigRadius,
  AgentConfigRadiusFields,
  AgentConfigRadiusToJSON,
  AgentConfigServerFields,
  AgentContainerConfiguration,
  AgentInstance,
  AzureRoleAssignment,
  CertPolicyRef,
  NamespaceKind,
  RadiusClientConfig,
  RadiusEapTls,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useCertPolicies } from "./CertPolicyRefTable";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import {
  NamespaceContext,
  NamespaceContextValue,
} from "./contexts/NamespaceContext";
import { XMarkIcon } from "@heroicons/react/24/outline";
import {
  RadiusConfigPatchContext,
  RadiusConfigPatchProvider,
  useRadiusConfigPatch,
} from "./contexts/RadiusConfigPatchContext";
import { RadiusConfigContainerForm } from "./forms/RadiusConfigContainerForm";

const wellKnownRoleDefinitionIds: Record<string, string> = {
  "7f951dda-4ed3-4680-a7ca-43fe172d538d": "AcrPull",
  "21090545-7ca7-4776-b22c-e363652d74d2": "Key Vault Reader",
  "4633458b-17de-408a-b874-0445c86b69e6": "Key Vault Secrets User",
};

export function useAzureRoleAssignmentsColumns(): ColumnsType<AzureRoleAssignment> {
  return useMemo(() => {
    return [
      {
        title: "Name",
        key: "name",
        render: (r: AzureRoleAssignment) => (
          <span className="font-mono">{r.name}</span>
        ),
      },
      {
        title: "Role definition id",
        key: "name",
        render: (r: AzureRoleAssignment) => {
          const parts = r.roleDefinitionId?.split("/");
          const defId = parts?.[parts.length - 1];
          return defId && wellKnownRoleDefinitionIds[defId] ? (
            wellKnownRoleDefinitionIds[defId]
          ) : (
            <span className="font-mono">{defId}</span>
          );
        },
      },
    ];
  }, []);
}

export type AgentServerConfigFormState = AgentContainerConfiguration;

function useCertPolicyOptions(
  certPolicies: CertPolicyRef[] | undefined
): DefaultOptionType[] | undefined {
  return certPolicies?.map((p) => ({
    label: p.displayName,
    value: p.id,
  }));
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
  }, [data]);
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
  }, [data]);
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
      <Form.Item<RadiusEapTlsFormState> name="certPolicyId" label="Server TLS certificate policy ID">
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
  }, [data]);
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
  const { data: radiusConfig, mutate } = patchSvc;

  const { data: keysData } = useRequest(
    () => {
      return api.listAgentAzureRoleAssignments({
        namespaceId,
        namespaceKind,
        configName: AgentConfigName.AgentConfigNameRadius,
      });
    },
    {
      ready: !!namespaceId && !!namespaceKind,
      refreshDeps: [namespaceId, namespaceKind],
    }
  );

  const roleAssignmentTableColumns = useAzureRoleAssignmentsColumns();

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
            <Card title="Azure role assignments">
              <JsonDataDisplay data={keysData} />
              <Table<AzureRoleAssignment>
                columns={roleAssignmentTableColumns}
                dataSource={keysData}
                rowKey={(r) => r.id ?? ""}
              />
            </Card>
          )}

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
