import { Button, Card, Table, Typography } from "antd";
import { useContext } from "react";
import { Link } from "../components/Link";
import {
  CertPolicyRef,
  ManagedAppRef,
  NamespaceKind,
  SecretPolicyRef,
} from "../generated";
import {
  useCertPolicies,
  usePolicyRefTableColumns,
  useSecretPolicies,
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
  const secretPolicyRoutePrefix = `/app/${NamespaceKind.NamespaceKindServicePrincipal}/${managedApp?.servicePrincipalId}/secret-policy/`;
  const secretPoliciesTableColumns = usePolicyRefTableColumns(
    secretPolicyRoutePrefix
  );

  const certPolicies = useCertPolicies();
  const secretPolicies = useSecretPolicies();
  return (
    <NamespaceConfigContextProvider ruleEntraClientCred>
      <section className="space-y-4">
        <Typography.Title level={2}>Service Principal</Typography.Title>
        <div>
          <pre>{managedApp?.servicePrincipalId}</pre>
        </div>
        <Card
          title="Certificate policies"
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
        <Card
          title="Secret policies"
          extra={
            <Link to={`${secretPolicyRoutePrefix}_create`}>
              Create secret policy
            </Link>
          }
        >
          <Table<SecretPolicyRef>
            columns={secretPoliciesTableColumns}
            dataSource={secretPolicies}
            rowKey={(r) => r.id}
          />
        </Card>
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
