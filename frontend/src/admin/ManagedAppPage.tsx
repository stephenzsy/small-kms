import { Button, Card, Table, Typography } from "antd";
import { useContext } from "react";
import { Link } from "../components/Link";
import {
  CertPolicyRef,
  ManagedAppRef,
  NamespaceKind
} from "../generated";
import {
  useCertPolicies,
  usePolicyRefTableColumns
} from "./CertPolicyRefTable";
import { ManagedAppContext } from "./contexts/ManagedAppContext";
import { NamespaceConfigContextProvider } from "./contexts/NamespaceConfigContextProvider";
import { NamespaceContext } from "./contexts/NamespaceContext";

function ServicePrincipalSection({
  managedApp,
}: {
  managedApp: ManagedAppRef | undefined;
}) {
  const certPolicyRoutePrefix = `/app/${NamespaceKind.NamespaceKindServicePrincipal}/${managedApp?.servicePrincipalId}/cert-policy/`;
  const columns = usePolicyRefTableColumns(certPolicyRoutePrefix);
  const certPolicies = useCertPolicies();
  return (
    <NamespaceConfigContextProvider ruleEntraClientCred>
      <section className="space-y-4">
        <Typography.Title level={2}>Service Principal</Typography.Title>
        <div>
          <pre>{managedApp?.servicePrincipalId}</pre>
        </div>
        <Card
          title="Certificate Policies"
          extra={
            <Link to={`${certPolicyRoutePrefix}_create`}>
              Create certificate policy
            </Link>
          }
        >
          <Table<CertPolicyRef>
            columns={columns}
            dataSource={certPolicies}
            rowKey={(r) => r.id}
          />
        </Card>
        <Card title="Provisioning">
          <Link to={`./provision-agent`}>Provision agent</Link>
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
          namespaceIdentifier: managedApp?.servicePrincipalId ?? "",
        }}
      >
        <ServicePrincipalSection managedApp={managedApp} />
      </NamespaceContext.Provider>
    </>
  );
}
