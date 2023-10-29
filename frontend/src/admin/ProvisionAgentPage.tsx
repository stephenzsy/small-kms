import { useContext, useMemo, useState } from "react";

import { Card, Form, Select, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import {
  AdminApi,
  AgentConfigName,
  AgentConfigurationAgentActiveHostBootstrapToJSON,
  AgentConfigurationAgentActiveServerToJSON,
  AgentConfigurationParameters,
  AgentConfigurationParametersFromJSON,
  CertPolicyRef,
  NamespaceKind1 as NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { useCertPolicies } from "./CertPolicyRefTable";
import { DefaultOptionType } from "antd/es/select";

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
        return JSON.stringify(
          AgentConfigurationAgentActiveServerToJSON({
            name: configName,
            authorizedCertificateTemplateId:
              "00000000-0000-0000-0000-000000000000",
            serverCertificateTemplateId: `cert-template:default-mtls`,
          }),
          undefined,
          2
        );
    }
    return "";
  }, [configName, nsId]);
}

export function AgentConfigurationForm({
  namespaceId,
  namespaceKind,
}: {
  namespaceId: string;
  namespaceKind: NamespaceKind;
}) {
  const adminApi = useAuthedClient(AdminApi);
  // const [selectedItem, setSelectedItem] = useState<SelectItem<AgentConfigName>>(
  //   selectOptions[0]
  // );

  // const currentConfigName = selectedItem.id;

  // const { data, loading, run } = useRequest(
  //   async (params?: AgentConfigurationParameters) => {
  //     if (params) {
  //       return await adminApi.putAgentConfiguration({
  //         configName: currentConfigName,
  //         namespaceKindLegacy: namespaceKind,
  //         namespaceId,
  //         agentConfigurationParameters: params,
  //       });
  //     }
  //     try {
  //       return await adminApi.getAgentConfiguration({
  //         namespaceKindLegacy: namespaceKind,
  //         namespaceId,
  //         configName: currentConfigName,
  //       });
  //     } catch (e) {
  //       return undefined;
  //     }
  //   },
  //   {
  //     refreshDeps: [namespaceId, namespaceKind, currentConfigName],
  //   }
  // );

  // const skeleton = useConfigurationSkeleton(
  //   currentConfigName,
  //   namespaceId,
  //   namespaceKind
  // );
  // const defaultValue = useMemo(() => {
  //   return loading
  //     ? ""
  //     : (data?.config
  //         ? JSON.stringify(
  //             AgentConfigurationParametersToJSON(data.config),
  //             undefined,
  //             2
  //           )
  //         : undefined) ?? skeleton;
  // }, [loading, data, skeleton]);
  const [configInput, setConfigInput] = useState<string>("");

  const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (!configInput.trim()) {
      return;
    }
    const parsed = JSON.parse(configInput);
    let typeParsed: AgentConfigurationParameters;

    typeParsed = AgentConfigurationParametersFromJSON(parsed);
    // run(typeParsed);
  };
  return null;
  /*<CardSection>
        Current configuration:
        {loading ? (
          <div>loading</div>
        ) : data ? (
          <div className="p-4">
            <pre>
              {JSON.stringify(AgentConfigurationToJSON(data), undefined, 2)}{" "}
            </pre>
          </div>
        ) : (
          <div>No configuration</div>
        )}
        </CardSection>
      <CardSection>
        <form className="space-y-4" onSubmit={onSubmit}>
          <Select
            label="Select Configuration Name"
            items={selectOptions}
            selected={selectedItem}
            setSelected={setSelectedItem}
          />
          <textarea
            className="w-full min-h-[400px]"
            value={configInput || defaultValue}
            onChange={(e) => {
              setConfigInput(e.target.value);
            }}
          />
          <Button type="submit" variant="primary">
            Update
          </Button>
        </form>
      </CardSection>*/
}

type AgentServerConfigFormState = {
  tlsCertPolicyId: string | undefined;
};

function useCertPolicyOptions(
  certPolicies: CertPolicyRef[] | undefined
): DefaultOptionType[] | undefined {
  return certPolicies?.map((p) => ({
    label: p.displayName,
    value: p.id,
  }));
}

function AgentConfigServerFormCard() {
  const [form] = useForm<AgentServerConfigFormState>();
  const certPolicies = useCertPolicies();
  const certPolicyOptions = useCertPolicyOptions(certPolicies);

  return (
    <Card title="Agent server configuration">
      <Form
        form={form}
        initialValues={{
          tlsCertPolicyId: undefined,
        }}
        layout="vertical"
      >
        <Form.Item<AgentServerConfigFormState>
          name="tlsCertPolicyId"
          label="Select server TLS certificate policy"
        >
          <Select options={certPolicyOptions} />
        </Form.Item>
      </Form>
    </Card>
  );
}

export default function ProvisionAgentPage() {
  const { managedApp } = useContext(ManagedAppContext);
  return (
    <>
      <Typography.Title>
        Provision agent: {managedApp?.displayName}
      </Typography.Title>

      <NamespaceContext.Provider
        value={{
          namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
          namespaceIdentifier: managedApp?.servicePrincipalId ?? "",
        }}
      >
        <AgentConfigServerFormCard />
      </NamespaceContext.Provider>
    </>
  );
}
