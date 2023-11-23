import { useMemoizedFn } from "ahooks";
import { Card, Table, TableColumnType } from "antd";
import { useMemo } from "react";
import { Link } from "../../components/Link";
import { ShortDate } from "../../components/ShortDate";
import { KeyRef, SecretRef } from "../../generated";

export function usePolicyItemRefTableColumns<T extends SecretRef | KeyRef>(
  onGetVewLink: (item: T) => string | undefined
) {
  const getLink = useMemoizedFn(onGetVewLink);
  return useMemo<TableColumnType<T>[]>(
    () => [
      {
        title: "ID",
        render: (r: T) => (
          <>
            <span className="font-mono">{r.id}</span>
          </>
        ),
      },
      {
        title: "Created",
        render: (r: T) => <ShortDate numericDate={r.iat} />,
      },
      {
        title: "Expires",
        render: (r: T) => <ShortDate numericDate={r.exp} />,
      },
      {
        render: (r: T) => {
          const link = getLink(r);
          if (link) {
            return <Link to={link}>View</Link>;
          }
          return null;
        },
      },
    ],
    []
  );
}

export function PolicyItemRefsTableCard<T extends SecretRef | KeyRef>({
  dataSource,
  title,
  onGetVewLink,
}: {
  dataSource: T[] | undefined;
  title: string;
  onGetVewLink: (item: T) => string | undefined;
}) {
  const columns = usePolicyItemRefTableColumns(onGetVewLink);
  return (
    <Card title={title}>
      <Table<T> columns={columns} dataSource={dataSource} rowKey="id" />
    </Card>
  );
}
