import { useRequest } from "ahooks";
import { Button, Card, Table, TableColumnType, Tag, Typography } from "antd";
import { useContext, useMemo } from "react";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import {
  AdminApi,
  CertificateRef,
  NamespaceKind,
  ResourceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { CertWebEnroll } from "./CertWebEnroll";
import {
  NamespaceConfigContext,
  NamespaceContext,
} from "./contexts/NamespaceContext";

const dateShortFormatter = new Intl.DateTimeFormat("en-US", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

function useCertTableColumns(activeIssuerCertificateId: string | undefined) {
  return useMemo<TableColumnType<CertificateRef>[]>(
    () => [
      {
        title: "Certificate ID",
        render: (r: CertificateRef) => (
          <>
            <span className="font-mono">{r.id}</span>
            {activeIssuerCertificateId === r.id && (
              <Tag className="ml-2" color="blue">
                Current issuer
              </Tag>
            )}
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
    ],
    [activeIssuerCertificateId]
  );
}

export default function CertPolicyPage() {
  const { certPolicyId: _certPolicyId } = useParams() as {
    certPolicyId: string;
  };
  const certPolicyId = _certPolicyId === "_create" ? "" : _certPolicyId;
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const { data } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.getCertPolicy({
          namespaceKind,
          resourceId: certPolicyId,
          namespaceId,
        });
      }
      return undefined;
    },
    {
      refreshDeps: [certPolicyId, namespaceId, namespaceKind],
    }
  );

  const { data: issuedCertificates } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.listCertificates({
          namespaceId,
          namespaceKind,
          policyId: certPolicyId,
        });
      }
      return undefined;
    },
    { refreshDeps: [namespaceId, certPolicyId] }
  );

  const {
    issuer: issuerRule,
    setIssuer: setIssuerRule,
    entraClientCred,
    setEntraClientCred,
  } = useContext(NamespaceConfigContext);

  const certListColumns = useCertTableColumns(issuerRule?.certificateId);

  return (
    <>
      <Typography.Title>
        Certificate Policy: {certPolicyId || "new policy"}
      </Typography.Title>
      <div className="font-mono">
        {namespaceKind}:{namespaceId}:{ResourceKind.ResourceKindCertPolicy}/
        {certPolicyId}
      </div>
      <Card title="Certificate list">
        <Table<CertificateRef>
          columns={certListColumns}
          dataSource={issuedCertificates}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card title="Manage certificates">
        <div className="space-y-4">
          {(namespaceKind === NamespaceKind.NamespaceKindRootCA ||
            namespaceKind === NamespaceKind.NamespaceKindIntermediateCA) &&
            (certPolicyId !== issuerRule?.policyId ? (
              <Button
                onClick={() => {
                  setIssuerRule({ policyId: certPolicyId });
                }}
              >
                Set as issuer policy
              </Button>
            ) : (
              <Tag color="blue">Current issuer policy</Tag>
            ))}
          {namespaceKind === NamespaceKind.NamespaceKindServicePrincipal && (
            <div className="flex gap-4 items-center">
              <Button
                onClick={() => {
                  setEntraClientCred({ policyId: certPolicyId });
                }}
              >
                Set as Microsoft Entra ID client credential policy
              </Button>
              {certPolicyId === entraClientCred?.policyId && (
                <Tag color="blue">
                  Current Microsoft Entra ID client credential policy
                </Tag>
              )}
            </div>
          )}
        </div>
      </Card>
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
      <Card title="Certificate web enrollment">
        <CertWebEnroll certPolicy={data} />
      </Card>
    </>
  );
}
