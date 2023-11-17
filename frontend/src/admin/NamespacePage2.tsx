import { useMemoizedFn, useRequest } from "ahooks";
import { Card, Typography } from "antd";
import { Link } from "../components/Link";
import { AdminApi, NamespaceProvider, Ref } from "../generated/apiv2";
import { useAuthedClientV2 } from "../utils/useCertsApi";
import { CertPolicyRefTable } from "./CertPolicyRefTable";
import { useNamespace } from "./contexts/NamespaceContextRouteProvider";
import { ResourceRefsTable } from "./tables/ResourceRefsTable";

export default function NamespacePage() {
  const { namespaceId, namespaceProvider } = useNamespace();

  const showCertPolicies = [
    NamespaceProvider.NamespaceProviderRootCA,
    NamespaceProvider.NamespaceProviderIntermediateCA,
    NamespaceProvider.NamespaceProviderServicePrincipal,
    NamespaceProvider.NamespaceProviderGroup,
  ].some((np) => np === namespaceProvider);
  const adminApi = useAuthedClientV2(AdminApi);
  const { data: certPolicies, loading: certPoliciesLoading } = useRequest(
    () => {
      return adminApi.listCertificatePolicies({
        namespaceId: namespaceId,
        namespaceProvider: namespaceProvider,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceProvider],
      ready: showCertPolicies,
    }
  );

  const renderActions = useMemoizedFn((ref: Ref) => {
    return (
      <div className="flex flex-row gap-2">
        <Link to={`./cert-policies/${ref.id}`}>View</Link>
      </div>
    );
  });

  return (
    <>
      <Typography.Title>Namespace: {namespaceId}</Typography.Title>
      <div className="font-mono">
        {namespaceProvider}:{namespaceId}
      </div>
      {showCertPolicies && (
        <Card
          title="Certificate Policies"
          extra={
            <Link to="./cert-policy/_create">Create certificate policy</Link>
          }
        >
          <ResourceRefsTable
            renderActions={renderActions}
            loading={certPoliciesLoading}
            dataSource={certPolicies}
          />
        </Card>
      )}
      {/* {namespaceKind === NamespaceKind.NamespaceKindUser && (
        <>
          <CertificatesTableCard />
          <Card title="Listed group memberships">
            <Table<ResourceReference>
              dataSource={groupMemberOf}
              columns={groupMemberOfColumns}
              rowKey="id"
            />
          </Card>
          <Card title="Sync group membership">
            <MemberOfGroupForm />
          </Card>
        </>
      )} */}
    </>
  );
}
