import { useContext, useEffect, useMemo, useState } from "react";

import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Divider,
  Form,
  Input,
  Select,
  Table,
  Typography,
} from "antd";
import { useForm } from "antd/es/form/Form";
import { DefaultOptionType } from "antd/es/select";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import {
  AdminApi,
  AgentConfigName,
  AgentConfigServerFields,
  AgentConfigServerToJSON,
  AgentConfigurationAgentActiveHostBootstrapToJSON,
  AgentConfigurationParameters,
  AgentConfigurationParametersFromJSON,
  AzureKeyvaultResourceCategory,
  AzureRoleAssignment,
  AzureRoleAssignmentFromJSON,
  CertPolicyRef,
  NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useCertPolicies } from "./CertPolicyRefTable";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import {
  NamespaceContext,
  NamespaceContextValue,
} from "./contexts/NamespaceContext";
import { ColumnsType } from "antd/es/table";
import { WellknownId } from "../constants";

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

function useConfigurationSkeleton(
  configName: AgentConfigName,
  nsId: string,
  nsKind: NamespaceKind
): string {
  return useMemo(() => {
    switch (configName) {
      case AgentConfigName.AgentConfigNameActiveHostBootstrap:
        return JSON.stringify(
          AgentConfigurationAgentActiveHostBootstrapToJSON({
            controllerContainer: {
              imageRefStr: "dockerimage:latest",
            },
            name: configName,
          }),
          undefined,
          2
        );
      case AgentConfigName.AgentConfigNameActiveServer:
      // return JSON.stringify(
      //   AgentConfigurationAgentActiveServerToJSON({
      //     name: configName,
      //     authorizedCertificateTemplateId:
      //       "00000000-0000-0000-0000-000000000000",
      //     serverCertificateTemplateId: `cert-template:default-mtls`,
      //   }),
      //   undefined,
      //   2
      // );
    }
    return "";
  }, [configName, nsId]);
}

const wellKnownRoleDefinitionIds: Record<string, string> = {
  "7f951dda-4ed3-4680-a7ca-43fe172d538d": "AcrPull",
  "21090545-7ca7-4776-b22c-e363652d74d2": "Key Vault Reader",
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

type AgentServerConfigFormState = Partial<AgentConfigServerFields> & {};

function useCertPolicyOptions(
  certPolicies: CertPolicyRef[] | undefined
): DefaultOptionType[] | undefined {
  return certPolicies?.map((p) => ({
    label: p.displayName,
    value: p.id,
  }));
}

function AgentConfigServerFormCard({
  isGlobalConfig,
}: {
  isGlobalConfig: boolean;
}) {
  const [form] = useForm<AgentServerConfigFormState>();
  const certPolicies = useCertPolicies();
  const certPolicyOptions = useCertPolicyOptions(certPolicies);
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { data: agentServerConfig, run } = useRequest(
    (params?: Partial<AgentConfigServerFields>) => {
      if (params) {
        return api.putAgentConfigServer({
          agentConfigServerFields: params as AgentConfigServerFields,
          namespaceKind,
          namespaceIdentifier,
        });
      }
      return api.getAgentConfigServer({
        namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier && !!namespaceKind,
    }
  );
  const jwtKeyCertPolicyId = agentServerConfig?.jwtKeyCertPolicyId;

  const { data: keysData } = useRequest(
    () => {
      return api.listAgentServerAzureRoleAssignments({
        namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      ready: !!namespaceIdentifier && !!jwtKeyCertPolicyId,
      refreshDeps: [namespaceIdentifier, jwtKeyCertPolicyId],
    }
  );

  useEffect(() => {
    if (agentServerConfig) {
      form.setFieldsValue(agentServerConfig);
    }
  }, [agentServerConfig]);

  const setCurrentBuildTag = useMemoizedFn(async () => {
    const tag = (await api.getDiagnostics()).serviceRuntime.buildId.split(
      "\\."
    )[0];
    const currentValue = form.getFieldValue("azureAcrImageRef");
    const currentPrfix = currentValue.split(":")[0];
    form.setFieldValue("azureAcrImageRef", `${currentPrfix}:${tag}`);
  });

  const roleAssignmentTableColumns = useAzureRoleAssignmentsColumns();
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
            name="tlsCertificatePolicyId"
            label="Select server TLS certificate policy"
            required
          >
            <Select options={certPolicyOptions} />
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
      <Divider />
      <JsonDataDisplay data={keysData} />
      <Table<AzureRoleAssignment>
        columns={roleAssignmentTableColumns}
        dataSource={keysData}
        rowKey={(r) => r.id ?? ""}
      />
    </Card>
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
            namespaceIdentifier: "default",
          }
        : {
            namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
            namespaceIdentifier: managedApp?.servicePrincipalId ?? "",
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
      </NamespaceContext.Provider>
    </>
  );
}
