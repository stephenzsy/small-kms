import { useRequest } from "ahooks";
import { Button } from "../components/Button";
import { Card, CardSection, CardTitle } from "../components/Card";
import Select, { SelectItem } from "../components/Select";
import {
  AgentConfigName,
  NamespaceKind,
  AdminApi,
  AgentConfiguration,
  AgentConfigurationParameters,
  AgentConfigurationToJSON,
  AgentConfigurationAgentActiveHostBootstrapToJSON,
  AgentConfigurationParametersFromJSON,
} from "../generated3";
import { useState, useMemo } from "react";
import { useAuthedClient } from "../utils/useCertsApi3";

const selectOptions: Array<SelectItem<AgentConfigName>> = [
  {
    id: AgentConfigName.AgentConfigNameActiveHostBootstrap,
    name: "Agent Active Host Bootstrap",
  },
];

function useConfigurationSkeleton(configName: AgentConfigName): string {
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
    }
    return "";
  }, [configName]);
}

export function AgentConfigurationForm({
  namespaceId,
  namespaceKind,
}: {
  namespaceId: string;
  namespaceKind: NamespaceKind;
}) {
  const adminApi = useAuthedClient(AdminApi);
  const { data, loading, run } = useRequest(
    (params?: AgentConfigurationParameters, configName?: AgentConfigName) => {
      if (params && configName) {
        return adminApi.putAgentConfiguration({
          configName,
          namespaceKind,
          namespaceId,
          agentConfigurationParameters: params,
        });
      }
      return adminApi.getAgentConfiguration({
        namespaceKind,
        namespaceId,
        configName: AgentConfigName.AgentConfigNameActiveHostBootstrap,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceKind],
    }
  );

  const [selectedItem, setSelectedItem] = useState<SelectItem<AgentConfigName>>(
    selectOptions[0]
  );

  const skeleton = useConfigurationSkeleton(selectedItem.id);
  const [configInput, setConfigInput] = useState<string>("");

  const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (!configInput.trim()) {
      return;
    }
    const parsed = JSON.parse(configInput);
    let typeParsed: AgentConfigurationParameters;
    switch (selectedItem.id) {
      case AgentConfigName.AgentConfigNameActiveHostBootstrap:
        typeParsed = AgentConfigurationParametersFromJSON(parsed);
        break;
    }
    run(typeParsed, selectedItem.id);
  };
  return (
    <Card>
      <CardTitle>Agent Configuration</CardTitle>
      <CardSection>
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
            value={configInput || skeleton}
            onChange={(e) => {
              setConfigInput(e.target.value);
            }}
          />
          <Button type="submit" variant="primary">
            Update
          </Button>
        </form>
      </CardSection>
    </Card>
  );
}
