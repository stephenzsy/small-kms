import { useMemoizedFn, useRequest } from "ahooks";
import { Card, Divider, Typography } from "antd";
import { useParams } from "react-router-dom";
import { useCertificatePolicies } from "../../admin/hooks/useCertificatePolicies";
import { ResourceRefsTable } from "../../admin/tables/ResourceRefsTable";
import { Link } from "../../components/Link";
import {
  AdminApi,
  NamespaceProvider,
  Ref as ResourceReference,
} from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

function useAgent(agentId: string | undefined) {
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    async () => {
      if (agentId) {
        return api.getAgent({
          id: agentId,
        });
      }
    },
    {
      refreshDeps: [agentId],
    }
  );
}

export default function AgentPage() {
  const { id } = useParams<{ id: string }>();
  const { data: agent } = useAgent(id);
  const spId = agent?.servicePrincipalId;
  const { data: certPolicies, loading: certPoliciesLoading } =
    useCertificatePolicies(
      NamespaceProvider.NamespaceProviderServicePrincipal,
      spId
    );
  const renderActions = useMemoizedFn((item: ResourceReference) => {
    return (
      <div className="flex flex-row gap-2">
        <Link to={`/service-principals/${spId}/cert-policies/${item.id}`}>
          View
        </Link>
      </div>
    );
  });

  return (
    <>
      <Typography.Title>Agent</Typography.Title>
      <Card title="Agent info">
        <dl className="dl">
          <div className="dl-item">
            <dt className="dt">ID (Client ID)</dt>
            <dd className="dd font-mono">{agent?.id}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Display name</dt>
            <dd className="dd">{agent?.displayName}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Application ID (Object ID)</dt>
            <dd className="dd font-mono">{agent?.applicationId}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Service Principal ID (Object ID)</dt>
            <dd className="dd font-mono">{agent?.servicePrincipalId}</dd>
          </div>
        </dl>
      </Card>
      <Divider />
      <Typography.Title level={2}>Service Principal</Typography.Title>
      <Card title="Certificate Policies">
        <ResourceRefsTable
          resourceRefs={certPolicies}
          loading={certPoliciesLoading}
          renderActions={renderActions}
        />
      </Card>
    </>
  );
}
