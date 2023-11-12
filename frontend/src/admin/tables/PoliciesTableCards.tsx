import { Button, Card, Table, type TableColumnType } from "antd";
import { AdminApi, KeyPolicy, KeyPolicyRef } from "../../generated";
import { useContext, useMemo } from "react";
import { Link } from "../../components/Link";
import {
  NamespaceConfigContext,
  NamespaceContext,
} from "../contexts/NamespaceContext";
import { useAuthedClient } from "../../utils/useCertsApi";
import { useRequest } from "ahooks";

export function usePolicyRefTableColumns<T extends KeyPolicyRef>(
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
      extra={
        <Link to={`${itemRoutePrefix}_create`}>
          Create key policy
        </Link>
      }
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
