import { Button, Card, Typography } from "antd";
import { useContext } from "react";
import { Link } from "../components/Link";
import {
  ManagedAppRef,
  NamespaceKind
} from "../generated";

import { ManagedAppContext } from "./contexts/ManagedAppContext";
import { NamespaceConfigContextProvider } from "./contexts/NamespaceConfigContextProvider";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { KeyPoliciesTabelCard } from "./tables/PoliciesTableCards";

function ServicePrincipalSection({
  managedApp,
}: {
  managedApp: ManagedAppRef | undefined;
}) {
  const keyPolicyRoutePrefix = `/app/${NamespaceKind.NamespaceKindServicePrincipal}/${managedApp?.servicePrincipalId}/key-policies/`;

  return (
    <NamespaceConfigContextProvider ruleEntraClientCred>
      <section className="space-y-4">
        <Typography.Title level={2}>Service Principal</Typography.Title>
        <div>
          <pre>{managedApp?.servicePrincipalId}</pre>
        </div>

        <KeyPoliciesTabelCard itemRoutePrefix={keyPolicyRoutePrefix} />

        <Card title="Provisioning">
          <div className="space-y-4 flex flex-col items-start">
            <Link to={`./provision-agent`}>Provision agent</Link>
            <Link to={`./radius-config`}>Configure radius</Link>
          </div>
        </Card>
      </section>
    </NamespaceConfigContextProvider>
  );
}

export default function ManagedAppPage({
  isSystemApp = false,
}: {
  isSystemApp?: boolean;
}) {
  const { managedApp, syncApp } = useContext(ManagedAppContext);

  return (
    <>
      <Typography.Title>
        {isSystemApp ? "System application: " : "Managed application: "}
        {managedApp?.displayName}
      </Typography.Title>
      <Card title="Sync">
        <Button type="primary" onClick={syncApp}>
          Sync
        </Button>
      </Card>
      <NamespaceContext.Provider
        value={{
          namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
          namespaceId: managedApp?.servicePrincipalId ?? "",
        }}
      >
        <ServicePrincipalSection managedApp={managedApp} />
      </NamespaceContext.Provider>
    </>
  );
}
