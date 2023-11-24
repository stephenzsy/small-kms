import { useRequest } from "ahooks";
import { Card, Typography } from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { AdminApi, ResourceKind } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";

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
        </>
      )}
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
    </>
  );
}
