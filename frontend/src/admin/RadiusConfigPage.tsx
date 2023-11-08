import { useContext, useEffect, useMemo } from "react";

import { useRequest } from "ahooks";
import {
  Button,
  Card,
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
  AgentConfigRadius,
  AgentConfigRadiusFields,
  AgentConfigRadiusToJSON,
  AgentConfigServerFields,
  AgentInstance,
  AzureRoleAssignment,
  CertPolicyRef,
  NamespaceKind,
  RadiusClientConfig,
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

const wellKnownRoleDefinitionIds: Record<string, string> = {
  "7f951dda-4ed3-4680-a7ca-43fe172d538d": "AcrPull",
  "21090545-7ca7-4776-b22c-e363652d74d2": "Key Vault Reader",
  "4633458b-17de-408a-b874-0445c86b69e6": "Key Vault Secrets User",
};

function useAzureRoleAssignmentsColumns(): ColumnsType<AzureRoleAssignment> {
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

type AgentServerConfigFormState = Pick<
  AgentConfigServerFields,
  "azureAcrImageRef"
>;

function useCertPolicyOptions(
  certPolicies: CertPolicyRef[] | undefined
): DefaultOptionType[] | undefined {
  return certPolicies?.map((p) => ({
    label: p.displayName,
    value: p.id,
  }));
}

function RadiusConfigGlobalFormCard({
  value,
  onUpdate,
}: {
  value: AgentConfigRadius | undefined;
  onUpdate?: (config: AgentConfigRadius) => void;
}) {
  const [form] = useForm<AgentServerConfigFormState>();
  // const certPolicies = useCertPolicies();
  // const certPolicyOptions = useCertPolicyOptions(certPolicies);
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { run } = useRequest(
    async (params: AgentConfigRadiusFields) => {
      const result = await api.putAgentConfigRadius({
        agentConfigRadiusFields: params,
        namespaceKind,
        namespaceId,
      });
      onUpdate?.(result);
    },
    {
      manual: true,
    }
  );

  // const { data: keysData } = useRequest(
  //   () => {
  //     return api.listAgentServerAzureRoleAssignments({
  //       namespaceId: namespaceIdentifier,
  //       namespaceKind,
  //     });
  //   },
  //   {
  //     ready: !!namespaceIdentifier && !!jwtKeyCertPolicyId,
  //     refreshDeps: [namespaceIdentifier, jwtKeyCertPolicyId],
  //   }
  // );

  useEffect(() => {
    if (value) {
      form.setFieldsValue({
        azureAcrImageRef: value.azureAcrImageRef,
      });
    }
  }, [value]);

  const roleAssignmentTableColumns = useAzureRoleAssignmentsColumns();
  return (
    <Card title="Global configuration">
      <Form form={form} layout="vertical" onFinish={run}>
        <Form.Item<AgentServerConfigFormState>
          name="azureAcrImageRef"
          label="Azure Container Registry image Reference"
          required
        >
          <Input placeholder="example.com/image:latest" />
        </Form.Item>

        {/* <Form.Item<AgentServerConfigFormState>
          name="tlsCertificatePolicyId"
          label="Select server TLS certificate policy"
          required
        >
          <Select options={certPolicyOptions} />
        </Form.Item> */}

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

function RadiusClientsForm() {
  const { run } = useRadiusConfigPatch();
  const [form] = useForm<Pick<AgentConfigRadiusFields, "clients">>();
  return (
    <Form form={form} layout="vertical" onFinish={(v) => run(v)}>
      <Form.List name={"clients"}>
        {(subFields, subOpt) => {
          return (
            <div className="flex flex-col gap-4">
              <div className="text-lg font-semibold">Secret bindings</div>
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
                    name={[subField.name, "secretRef"]}
                    className="flex-auto"
                    label="Secret ID"
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
    key: "clients",
    label: "Clients",
    children: <RadiusClientsForm />,
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
        return api.putAgentConfigRadius({
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
          {isGlobalConfig ? (
            <RadiusConfigGlobalFormCard
              value={radiusConfig}
              onUpdate={mutate}
            />
          ) : (
            <Card title="RADIUS configuration">
              <Collapse items={collapseItems} />
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
