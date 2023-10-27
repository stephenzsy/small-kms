import { Button, Table, Tag, type TableColumnType } from "antd";
import { AdminApi, CertPolicyRef, NamespaceKind } from "../generated";
import { useMemo, useContext } from "react";
import { Link } from "../components/Link";
import { NamespaceContext } from "./NamespaceContext";
import { useAuthedClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";
import { CertificateIssuerContext } from "./CertIssuerContext";

function useColumns(
  routePrefix: string,
  activeIssuerPolicyId: string | undefined
) {
  return useMemo<TableColumnType<CertPolicyRef>[]>(
    () => [
      {
        title: "App ID",
        render: (r: CertPolicyRef) => (
          <>
            <span className="font-mono">{r.id}</span>
            {r.id === activeIssuerPolicyId && (
              <Tag className="ml-2" color="blue">Current issuer</Tag>
            )}
          </>
        ),
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
    [routePrefix, activeIssuerPolicyId]
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

  const { rule: issuerRule } = useContext(CertificateIssuerContext);

  const columns = useColumns(routePrefix, issuerRule?.policyId);
  return (
    <Table<CertPolicyRef>
      columns={columns}
      dataSource={certPolicies}
      rowKey={(r) => r.id}
    />
  );
}
