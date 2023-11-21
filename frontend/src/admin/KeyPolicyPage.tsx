import { useRequest } from "ahooks";
import { Button, Card, Typography } from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import {
  AdminApi,
  ResourceKind
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { PolicyItemRefsTableCard } from "./tables/PolicyItemRefsTableCard";

function GenerateKeyControl({
  policyId,
  onComplete,
}: {
  policyId: string;
  onComplete?: () => void;
}) {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const { run, loading } = useRequest(
    async () => {
      await adminApi.generateKey({
        resourceId: policyId,
        namespaceId: namespaceId,
        namespaceKind,
      });
      onComplete?.();
    },
    { manual: true }
  );

  return (
    <div className="flex gap-8 items-center">
      <Button
        loading={loading}
        type="primary"
        onClick={() => {
          run();
        }}
      >
        Generate cloud key
      </Button>
    </div>
  );
}

export default function KeyPolicyPage() {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const _policyId = useParams().policyId;
  const policyId = _policyId === "_create" ? "" : _policyId;

  const api = useAuthedClient(AdminApi);
  const {
    data,
    run: refresh,
    mutate,
  } = useRequest(
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

  const { data: issuedKeys, run: refreshKeys } = useRequest(
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
          <Card title="Manage key">
            {policyId && (
              <GenerateKeyControl
                policyId={policyId}
                onComplete={refreshKeys}
              />
            )}
          </Card>
        </>
      )}
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
    </>
  );
}
