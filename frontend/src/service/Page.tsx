import { useRequest } from "ahooks";
import { useEffect, useState, type PropsWithChildren } from "react";
import { Link } from "react-router-dom";
import { InputField } from "../admin/InputField";
import { Button } from "../components/Button";
import { Card, CardSection, CardTitle } from "../components/Card";
import {
  AdminApi,
  NamespaceKind,
  PatchServiceConfigConfigPathEnum,
  ServiceConfig,
  ServiceConfigAcrInfoToJSON,
  ServiceConfigAppRoleIdsToJSON,
  ServiceConfigToJSON,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

function PatchServiceConfigForm(
  props: PropsWithChildren<{ onSubmit: () => void }>
) {
  return (
    <form
      className="space-y-4 p-4"
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        props.onSubmit();
      }}
    >
      {props.children}
      <Button type="submit" variant="primary">
        Update
      </Button>
    </form>
  );
}

export default function ServicePage() {
  const adminApi = useAuthedClient(AdminApi);
  const { data: serviceConfig, run: patchServiceConfig } = useRequest(
    (
      patchKey?: PatchServiceConfigConfigPathEnum,
      data?: ServiceConfig[PatchServiceConfigConfigPathEnum]
    ) => {
      if (patchKey && data) {
        return adminApi.patchServiceConfig({
          configPath: patchKey,
          body: data,
        });
      }
      return adminApi.getServiceConfig();
    }
  );

  const [azureSubscriptionId, setAzureSubscriptionId] = useState<string>("");
  const [keyVaultResourceId, setKeyVaultResourceId] = useState<string>("");

  const [appRoleIdAppAdmin, setAppRoleIdAppAdmin] = useState<string>("");
  const [appRoleIdAgentActiveHost, setAppRoleIdAgentActiveHost] =
    useState<string>("");

  const [acrLoginServer, setAcrLoginServer] = useState<string>("");
  const [acrName, setAcrName] = useState<string>("");
  const [acrResourceId, setAcrResourceId] = useState<string>("");

  useEffect(() => {
    if (serviceConfig) {
      setAzureSubscriptionId(serviceConfig.azureSubscriptionId);
      setKeyVaultResourceId(serviceConfig.keyvaultArmResourceId);
      setAppRoleIdAppAdmin(serviceConfig.appRoleIds.appAdmin);
      setAppRoleIdAgentActiveHost(serviceConfig.appRoleIds.agentActiveHost);
      setAcrLoginServer(serviceConfig.azureContainerRegistry.loginServer);
      setAcrName(serviceConfig.azureContainerRegistry.name);
      setAcrResourceId(serviceConfig.azureContainerRegistry.armResourceId);
    }
  }, [serviceConfig]);

  return (
    <>
      <Card>
        <CardTitle>Service configuration</CardTitle>
        <CardSection>
          <Link
            to={`/admin/${NamespaceKind.NamespaceKindSystem}/agent-push/certificate-templates/default-mtls`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            Service mTLS Certificate template
          </Link>
        </CardSection>
        <CardSection>
          Current configuration:
          <pre>
            {JSON.stringify(ServiceConfigToJSON(serviceConfig), null, 2)}
          </pre>
        </CardSection>
        <CardSection>
          <PatchServiceConfigForm
            onSubmit={() => {
              patchServiceConfig(
                PatchServiceConfigConfigPathEnum.ServiceConfigPathAzureSubscriptionId,
                azureSubscriptionId
              );
            }}
          >
            <InputField
              labelContent="Azure subscription ID"
              value={azureSubscriptionId}
              onChange={setAzureSubscriptionId}
            />
          </PatchServiceConfigForm>
        </CardSection>
        <CardSection>
          <PatchServiceConfigForm
            onSubmit={() => {
              patchServiceConfig(
                PatchServiceConfigConfigPathEnum.ServiceConfigPathKeyvaultArmResourceId,
                keyVaultResourceId
              );
            }}
          >
            <InputField
              labelContent="Key vault resource ID"
              value={keyVaultResourceId}
              onChange={setKeyVaultResourceId}
            />
          </PatchServiceConfigForm>
        </CardSection>
        <CardSection>
          <PatchServiceConfigForm
            onSubmit={() => {
              patchServiceConfig(
                PatchServiceConfigConfigPathEnum.ServiceConfigPathAppRoleIds,
                ServiceConfigAppRoleIdsToJSON({
                  appAdmin: appRoleIdAppAdmin,
                  agentActiveHost: appRoleIdAgentActiveHost,
                })
              );
            }}
          >
            <h3 className="text-lg font-semibold">App role IDs</h3>
            <InputField
              labelContent="App.Admin"
              value={appRoleIdAppAdmin}
              onChange={setAppRoleIdAppAdmin}
            />
            <InputField
              labelContent="Agent.ActiveHost"
              value={appRoleIdAgentActiveHost}
              onChange={setAppRoleIdAgentActiveHost}
            />
          </PatchServiceConfigForm>
        </CardSection>
        <CardSection>
          <PatchServiceConfigForm
            onSubmit={() => {
              patchServiceConfig(
                PatchServiceConfigConfigPathEnum.ServiceConfigPathAzureContainerRegistry,
                ServiceConfigAcrInfoToJSON({
                  loginServer: acrLoginServer,
                  name: acrName,
                  armResourceId: acrResourceId,
                })
              );
            }}
          >
            <h3 className="text-lg font-semibold">Azure container registry</h3>
            <InputField
              labelContent="Name"
              value={acrName}
              onChange={setAcrName}
            />
            <InputField
              labelContent="Login server"
              value={acrLoginServer}
              onChange={setAcrLoginServer}
            />
            <InputField
              labelContent="Resource ID"
              value={acrResourceId}
              onChange={setAcrResourceId}
            />
          </PatchServiceConfigForm>
        </CardSection>
      </Card>
    </>
  );
}
