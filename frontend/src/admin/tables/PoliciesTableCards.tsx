import { useRequest } from "ahooks";
import { Card, Table, type TableColumnType } from "antd";
import { useContext, useMemo } from "react";
import { Link } from "../../components/Link";
import { AdminApi, KeyPolicyRef } from "../../generated";
import { useAuthedClient } from "../../utils/useCertsApi";
import { NamespaceContext } from "../contexts/NamespaceContext";

function usePolicyRefTableColumns<T extends KeyPolicyRef>(
  routePrefix: string,
  onRenderTags?: (r: T) => React.ReactNode
) {
  return useMemo<TableColumnType<T>[]>(
    () => [
      {
        title: "ID",
        render: (r: T) => (
          <>
            <span className="font-mono">{r.id}</span>
            {onRenderTags?.(r)}
          </>
        ),
      },
      {
        title: "Display name",
        render: (r: T) => r.displayName,
      },

      {
        title: "Actions",
        render: (r: T) => (
          <>
            <Link to={`${routePrefix}${r.id}`}>View</Link>
          </>
        ),
      },
    ],
    [routePrefix, onRenderTags]
  );
}

export function KeyPoliciesTabelCard({
  itemRoutePrefix,
}: {
  itemRoutePrefix: string;
}) {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data, loading } = useRequest(
    () => {
      return adminApi.listKeyPolicies({
        namespaceId: namespaceId,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceKind],
      ready: !!namespaceId && !!namespaceKind,
    }
  );

  const columns = usePolicyRefTableColumns(itemRoutePrefix);

  return (
    <Card
      title="Key policies"
      extra={<Link to={`${itemRoutePrefix}_create`}>Create key policy</Link>}
    >
      <Table<KeyPolicyRef>
        columns={columns}
        dataSource={data}
        loading={loading}
        rowKey={(r) => r.id}
      />
    </Card>
  );
}
