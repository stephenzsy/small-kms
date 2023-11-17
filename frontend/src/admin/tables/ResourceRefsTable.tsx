import { Table, TableColumnsType } from "antd";
import { useMemo } from "react";
import { Ref as ResourceRef } from "../../generated/apiv2";

function useResourceRefsTableColumns<T extends ResourceRef>(
  renderActions?: (r: T) => React.ReactNode,
  extraColumns?: TableColumnsType<T>
) {
  return useMemo((): TableColumnsType<T> => {
    return [
      {
        title: "ID",
        dataIndex: "id",
        key: "id",
        className: "font-mono",
      },
      {
        title: "Display Name",
        dataIndex: "displayName",
        key: "displayName",
      },
      ...(extraColumns ?? []),
      ...(renderActions
        ? [
            {
              title: "Actions",
              key: "actions",
              render: renderActions,
            },
          ]
        : []),
    ];
  }, [extraColumns, renderActions]);
}

export function ResourceRefsTable<T extends ResourceRef>({
  dataSource,
  loading,
  renderActions,
  extraColumns,
}: {
  dataSource: T[] | undefined;
  loading?: boolean;
  renderActions?: (r: T) => React.ReactNode;
  extraColumns?: TableColumnsType<T>;
}) {
  const columns = useResourceRefsTableColumns<T>(renderActions, extraColumns);
  return (
    <Table<T>
      dataSource={dataSource}
      loading={loading}
      columns={columns}
      rowKey="id"
    />
  );
}
