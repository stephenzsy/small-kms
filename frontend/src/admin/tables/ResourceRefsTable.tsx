import { Table, TableColumnsType } from "antd";
import { Ref as ResourceRef } from "../../generated/apiv2";
import { useMemo } from "react";

function useResourceRefsTableColumns<T extends ResourceRef>(
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
    ];
  }, [extraColumns]);
}

export function ResourceRefsTable<T extends ResourceRef>({
  resourceRefs,
  loading,
  extraColumns,
}: {
  resourceRefs?: T[];
  loading?: boolean;
  extraColumns?: TableColumnsType<T>;
}) {
  const columns = useResourceRefsTableColumns<T>(extraColumns);
  return (
    <Table<T> dataSource={resourceRefs} loading={loading} columns={columns} />
  );
}
