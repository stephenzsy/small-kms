import { Card, Typography } from "antd";
import { NamespaceContext, NamespaceContextProvider } from "./NamespaceContext";
import { AdminApi, NamespaceKind } from "../generated";
import { useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi";
import { useParams } from "react-router-dom";
import { Link } from "../components/Link";
import { CertPolicyRefTable } from "./CertPolicyRefTable";

export default function ManagedAppPage() {
  const { appId } = useParams() as { appId: string };
  const adminApi = useAuthedClient(AdminApi);
  const { data: managedApp } = useRequest(() => {
    return adminApi.getManagedApp({ managedAppId: appId });
  }, {});

  return (
    <>
      <Typography.Title>
        Managed application: {managedApp?.displayName}
      </Typography.Title>
      <NamespaceContext.Provider
        value={{
          namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
          namespaceIdentifier: managedApp?.servicePrincipalId ?? "",
        }}
      >
        <section>
          <Typography.Title level={2}>Service Principal</Typography.Title>
          <div>
            <pre>{managedApp?.servicePrincipalId}</pre>
          </div>
          <Card
            title="Certificate Policies"
            extra={
              <Link
                to={`/app/${NamespaceKind.NamespaceKindServicePrincipal}/${managedApp?.servicePrincipalId}/cert-policy/_create`}
              >
                Create certificate policy
              </Link>
            }
          >
            <CertPolicyRefTable />
          </Card>
        </section>
      </NamespaceContext.Provider>
    </>
  );
}
