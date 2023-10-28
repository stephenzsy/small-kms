import { useRequest } from "ahooks";
import { Button, Card, Form, Table, Typography } from "antd";
import { useParams } from "react-router-dom";
import { Link } from "../components/Link";
import {
  AdminApi,
  CertPolicyRef,
  ManagedAppRef,
  NamespaceKind,
  SystemAppName,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  CertPolicyRefTable,
  useCertPolicies,
  usePolicyRefTableColumns,
} from "./CertPolicyRefTable";
import { NamespaceConfigContextProvider } from "./contexts/NamespaceConfigContextProvider";
import { NamespaceContext } from "./contexts/NamespaceContext";
import Select, { DefaultOptionType } from "antd/es/select";
import { useForm } from "antd/es/form/Form";

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

function AgentServerConfigForm({
  certPolicies,
}: {
  certPolicies: CertPolicyRef[] | undefined;
}) {
  const certPolicyOptions = useCertPolicyOptions(certPolicies);
  const [form] = useForm<AgentServerConfigFormState>();
  return (
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
  );
}

function ServicePrincipalSection({
  managedApp,
}: {
  managedApp: ManagedAppRef | undefined;
}) {
  const routePrefix = `/app/${NamespaceKind.NamespaceKindServicePrincipal}/${managedApp?.servicePrincipalId}/cert-policy/`;
  const columns = usePolicyRefTableColumns(routePrefix);
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
            <Link to={`${routePrefix}_create`}>Create certificate policy</Link>
          }
        >
          <Table<CertPolicyRef>
            columns={columns}
            dataSource={certPolicies}
            rowKey={(r) => r.id}
          />
        </Card>
        <Card title="Agent server configuration">
          <AgentServerConfigForm certPolicies={certPolicies} />
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
  const { appId } = useParams() as { appId: string };
  const adminApi = useAuthedClient(AdminApi);
  const { data: managedApp, run } = useRequest(
    (sync?: boolean) => {
      if (isSystemApp) {
        if (sync) {
          return adminApi.syncSystemApp({
            systemAppName: appId as SystemAppName,
          });
        }
        return adminApi.getSystemApp({ systemAppName: appId as SystemAppName });
      }
      if (sync) {
        return adminApi.syncManagedApp({ managedAppId: appId });
      }
      return adminApi.getManagedApp({ managedAppId: appId });
    },
    {
      refreshDeps: [appId, isSystemApp],
    }
  );

  return (
    <>
      <Typography.Title>
        Managed application: {managedApp?.displayName}
      </Typography.Title>
      <Card title="Sync">
        <Button
          type="primary"
          onClick={() => {
            run(true);
          }}
        >
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
