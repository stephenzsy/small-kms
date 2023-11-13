import { Card, Table, TableColumnType, Tag } from "antd";
import { AdminApi, CertificateRef } from "../../generated";
import { useContext, useMemo } from "react";
import { Link } from "../../components/Link";
import { useRequest } from "ahooks";
import { NamespaceContext } from "../contexts/NamespaceContext";
import { useAuthedClient } from "../../utils/useCertsApi";
import { dateShortFormatter } from "../../utils/datetimeUtils";


type CertificateActionsProps = {
  certRef: CertificateRef;
  onSetIssuerPolicy?: (policyId: string) => void;
  certPolicyId?: string;
};

function CertificateActions({
  certRef,
  certPolicyId,
  onSetIssuerPolicy,
}: CertificateActionsProps) {
  return (
    <div className="flex items-center gap-2">
      <Link to={`/my/certs/${certRef.id}`}>View</Link>
    </div>
  );
}

function useCertTableColumns(activeIssuerCertificateId: string | undefined) {
  return useMemo<TableColumnType<CertificateRef>[]>(
    () => [
      {
        title: "Certificate ID",
        render: (r: CertificateRef) => (
          <>
            <span className="font-mono">{r.id}</span>
          </>
        ),
      },
      {
        title: "Thumbprint (SHA-1)",
        render: (r: CertificateRef) => {
          return <span className="font-mono">{r.thumbprint}</span>;
        },
      },
      {
        title: "Expires",
        render: (r: CertificateRef) => {
          return (
            <span className="font-mono">
              {r.attributes.exp &&
                dateShortFormatter.format(new Date(r.attributes.exp * 1000))}
            </span>
          );
        },
      },
      {
        title: "Status",
        render: (r: CertificateRef) => {
          if (r.deleted) {
            return <Tag color="red">Deleted</Tag>;
          } else if (!r.thumbprint) {
            return <Tag color="yellow">Pending</Tag>;
          }
          return <Tag color="green">Issued</Tag>;
        },
      },
      {
        title: "Actions",
        render: (r) => <CertificateActions certRef={r} />,
      },
    ],
    [activeIssuerCertificateId]
  );
}

export function CertificatesTableCard({
  certPolicyId,
}: {
  certPolicyId?: string;
}) {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const api = useAuthedClient(AdminApi);
  const { data: issuedCertificates, refresh: refreshCertificates } = useRequest(
    async () => {
      return await api.listCertificates({
        namespaceId,
        namespaceKind,
        policyId: certPolicyId,
      });
    },
    { refreshDeps: [namespaceId, certPolicyId] }
  );
  const columns = useCertTableColumns(undefined);
  return (
    <Card title="Certificate list">
      <Table<CertificateRef>
        columns={columns}
        dataSource={issuedCertificates}
        rowKey="id"
      />
    </Card>
  );
}
