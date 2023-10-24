import { useRequest } from "ahooks";
import { Card, Table, TableColumnType } from "antd";
import { useContext, useMemo } from "react";
import { Link } from "../components/Link";
import { AdminApi, CertPolicyRef } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./NamespaceContext2";

function useColumns() {
  return useMemo<TableColumnType<CertPolicyRef>[]>(
    () => [
      {
        title: "App ID",
        render: (r: CertPolicyRef) => (
          <span className="font-mono">{r.resourceIdentifier}</span>
        ),
      },
      {
        title: "Display name",
        render: (r: CertPolicyRef) => r.displayName,
      },

      {
        title: "Actions",
        render: (r: CertPolicyRef) => (
          <Link to={`./cert-policy/${r.resourceIdentifier}`}>View</Link>
        ),
      },
    ],
    []
  );
}

export default function NamespacePage() {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: certPolicies } = useRequest(
    () => {
      return adminApi.listCertPolicies({
        namespaceIdentifier: namespaceId,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceKind],
    }
  );
  const columns = useColumns();

  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{namespaceKind}</div>
      <Card
        title="Certificate Policies"
        extra={
          <Link to="./cert-policy/_create">Create certificate policy</Link>
        }
      >
        <Table<CertPolicyRef>
          columns={columns}
          dataSource={certPolicies}
          rowKey={(r) => r.resourceIdentifier}
        />
      </Card>
    </>
  );
}
