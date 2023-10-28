import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { Card, CardSection, CardTitle } from "../components/Card";
import {
  AdminApi,
  AgentConfigName,
  AgentConfigurationAgentActiveHostBootstrapToJSON,
  AgentConfigurationAgentActiveServerToJSON,
  AgentConfigurationParameters,
  AgentConfigurationParametersFromJSON,
  NamespaceKind1 as NamespaceKind
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

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
  return (
    <Card>
      <CardTitle>Agent Configuration</CardTitle>
      <CardSection>
        <Link
          to={`/admin/${namespaceKind}/${namespaceId}/agent`}
          className="text-indigo-600 hover:text-indigo-900"
        >
          Go to agent dashboard
        </Link>
      </CardSection>
      {/*<CardSection>
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
      </CardSection>*/}
    </Card>
  );
}
