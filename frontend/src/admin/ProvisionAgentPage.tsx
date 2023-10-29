import { useContext, useEffect, useMemo, useState } from "react";

import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Select, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import { DefaultOptionType } from "antd/es/select";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import {
  AdminApi,
  AgentConfigName,
  AgentConfigServerParameters,
  AgentConfigServerToJSON,
  AgentConfigurationAgentActiveHostBootstrapToJSON,
  AgentConfigurationParameters,
  AgentConfigurationParametersFromJSON,
  CertPolicyRef,
  NamespaceKind1 as NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useCertPolicies } from "./CertPolicyRefTable";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import { NamespaceContext } from "./contexts/NamespaceContext";

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

type AgentServerConfigFormState = AgentConfigServerParameters & {};

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
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { data, run } = useRequest(
    (params?: AgentConfigServerParameters) => {
      if (params) {
        return api.putAgentConfigServer({
          agentConfigServerParameters: params,
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

  useEffect(() => {
    if (data) {
      form.setFieldsValue(data);
    }
  }, [data])

  return (
    <Card title="Agent server configuration">
      <div className="mb-6">
        Current configuration:
        <JsonDataDisplay data={data} toJson={AgentConfigServerToJSON} />
      </div>
      <Form
        form={form}
        layout="vertical"
        onFinish={(values) => {
          run(values);
        }}
      >
        <Form.Item<AgentServerConfigFormState>
          name="tlsCertificatePolicyId"
          label="Select server TLS certificate policy"
          required
        >
          <Select options={certPolicyOptions} />
        </Form.Item>
        <Form.Item<AgentServerConfigFormState>
          name="jwtKeyCertPolicyId"
          label="Json web token key certificate policy"
          required
        >
          <Input />
        </Form.Item>
        <Form.Item>
          <Button htmlType="submit" type="primary">
            Submit
          </Button>
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
