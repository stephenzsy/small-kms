import { Button, Table, type TableColumnType } from "antd";
import { AdminApi, CertPolicyRef, NamespaceKind } from "../generated";
import { useMemo, useContext } from "react";
import { Link } from "../components/Link";
import { NamespaceContext } from "./NamespaceContext";
import { useAuthedClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";

function useColumns(routePrefix: string) {
  return useMemo<TableColumnType<CertPolicyRef>[]>(
    () => [
      {
        title: "App ID",
        render: (r: CertPolicyRef) => <span className="font-mono">{r.id}</span>,
      },
      {
        title: "Display name",
        render: (r: CertPolicyRef) => r.displayName,
      },

      {
        title: "Actions",
        render: (r: CertPolicyRef) => (
          <>
            <Link to={`${routePrefix}${r.id}`}>View</Link>
          </>
        ),
      },
    ],
    [routePrefix]
  );
}

export function CertPolicyRefTable({ routePrefix }: { routePrefix: string }) {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: certPolicies } = useRequest(
    () => {
      return adminApi.listCertPolicies({
        namespaceIdentifier,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier,
    }
  );

  const { data } = useRequest(
    () => {
      return adminApi.getCertificateRuleIssuer({
        namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready:
        namespaceKind === NamespaceKind.NamespaceKindRootCA ||
        namespaceKind === NamespaceKind.NamespaceKindIntermediateCA,
    }
  );
  const columns = useColumns(routePrefix);
  return (
    <Table<CertPolicyRef>
      columns={columns}
      dataSource={certPolicies}
      rowKey={(r) => r.id}
    />
  );
}
