import { useRequest } from "ahooks";
import { Card, Typography } from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { AdminApi, ResourceKind } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { PolicyItemRefsTableCard } from "./tables/PolicyItemRefsTableCard";

export default function KeyPolicyPage() {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const _policyId = useParams().policyId;
  const policyId = _policyId === "_create" ? "" : _policyId;

  const api = useAuthedClient(AdminApi);
  const { data } = useRequest(
    async () => {
      return await api.getKeyPolicy({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        resourceId: policyId!,
      });
    },
    {
      refreshDeps: [policyId, namespaceIdentifier, namespaceKind],
      ready: !!policyId,
    }
  );

  const { data: issuedKeys } = useRequest(
    async () => {
      return await api.listKeys({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        policyId: policyId!,
      });
    },
    {
      refreshDeps: [policyId, namespaceIdentifier, namespaceKind],
      ready: !!policyId,
    }
  );

  return (
    <>
      <Typography.Title>
        Key Policy: {policyId || "new policy"}
      </Typography.Title>
      {policyId && (
        <>
          <div className="font-mono">
            {namespaceKind}:{namespaceIdentifier}:
            {ResourceKind.ResourceKindKeyPolicy}/{policyId}
          </div>
          <PolicyItemRefsTableCard
            title="Key list"
            dataSource={issuedKeys}
            onGetVewLink={(r) => `../keys/${r.id}`}
          />
        </>
      )}
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
    </>
  );
}
