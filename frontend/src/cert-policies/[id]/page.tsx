import {
  Button,
  Card,
  Form,
  Input,
  TableColumnType,
  TableColumnsType,
  Tag,
  Typography,
} from "antd";
import { useNamespace } from "../../admin/contexts/NamespaceContextRouteProvider";
import { AdminApi, CertificateStatus } from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useRequest } from "ahooks";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { CertPolicyForm } from "../../admin/forms/CertPolicyForm";
import { useContext } from "react";
import { DrawerContext } from "../../admin/contexts/DrawerContext";
import { ResourceRefsTable } from "../../admin/tables/ResourceRefsTable";
import { CertificateRef } from "../../generated/apiv2";
import { dateShortFormatter } from "../../utils/datetimeUtils";
import { Link } from "../../components/Link";

function useCertificatePolicy(id: string | undefined) {
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    async () => {
      if (id) {
        return api.getCertificatePolicy({
          id: id,
          namespaceId: namespaceId,
          namespaceProvider: namespaceProvider,
        });
      }
    },
    {
      refreshDeps: [id, namespaceId, namespaceProvider],
    }
  );
}

function useCertificateTableColumns(): TableColumnsType<CertificateRef> {
  return [
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        return (
          <Tag
            className="capitalize"
            color={
              status === CertificateStatus.CertificateStatusIssued
                ? "green"
                : status === CertificateStatus.CertificateStatusPending
                ? "orange"
                : status === CertificateStatus.CertificateStatusDeactivated
                ? "red"
                : undefined
            }
          >
            {status}
          </Tag>
        );
      },
    },
    {
      title: "Thumbprint SHA-1",
      dataIndex: "thumbprint",
      key: "thumbprint",
      render: (tp?: string) => {
        return <span className="font-mono text-xs uppercase">{tp}</span>;
      },
    },
    {
      title: "Date Issued",
      dataIndex: "iat",
      key: "iat",
      render: (tsNum?: number) => {
        if (tsNum) {
          const ts = new Date(tsNum * 1000);
          return (
            <time dateTime={ts.toISOString()}>
              {dateShortFormatter.format(ts)}
            </time>
          );
        }
      },
    },
    {
      title: "Date Expires",
      dataIndex: "exp",
      key: "exp",
      render: (tsNum?: number) => {
        if (tsNum) {
          const ts = new Date(tsNum * 1000);
          return (
            <time dateTime={ts.toISOString()}>
              {dateShortFormatter.format(ts)}
            </time>
          );
        }
      },
    },
  ];
}

export default function CertPolicyPage() {
  const { id } = useParams<{ id: string }>();
  const { data: certPolicy, mutate } = useCertificatePolicy(id);
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const {
    data: certs,
    refresh: refreshCerts,
    loading: certsLoading,
  } = useRequest(
    async () => {
      return await api.listCertificates({
        namespaceId: namespaceId,
        namespaceProvider: namespaceProvider,
        policyId: id,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceProvider, id],
    }
  );
  const { run: generateCertificate, loading: generateCertificateLoading } =
    useRequest(
      async () => {
        if (id && certPolicy) {
          return await api.generateCertificate({
            id: id,
            namespaceId: namespaceId,
            namespaceProvider: namespaceProvider,
          });
        }
        refreshCerts();
      },
      { manual: true }
    );

  const certColumns = useCertificateTableColumns();
  const viewCert = (cert: CertificateRef) => {
    return (
      <div>
        <Link to={`../certificates/${cert.id}`}>View</Link>
      </div>
    );
  };
  return (
    <>
      <Typography.Title>Certificate Policy</Typography.Title>
      <div className="font-mono">
        {namespaceProvider}:{namespaceId}:cert-policy/{id}
      </div>
      <Card
        title="Certificates"
        extra={
          <>
            {certPolicy?.allowGenerate && (
              <Button
                type="primary"
                onClick={generateCertificate}
                loading={generateCertificateLoading}
              >
                Generate Certificate
              </Button>
            )}
          </>
        }
      >
        <ResourceRefsTable<CertificateRef>
          loading={certsLoading}
          dataSource={certs}
          extraColumns={certColumns}
          noDisplayName
          renderActions={viewCert}
        />
      </Card>
      <Card
        title="Certificate Policy"
        extra={
          <Button
            type="link"
            onClick={() =>
              openDrawer(
                <div className="space-y-4">
                  <label>
                    <span className="text-sm mb-2">Policy ID:</span>
                    <Input
                      readOnly
                      className="font-mono"
                      value={`${namespaceProvider}:${namespaceId}:cert-policy/${certPolicy?.id}`}
                    />
                  </label>
                  <JsonDataDisplay data={certPolicy} />
                </div>,
                {
                  title: "Certificate Policy",
                  size: "large",
                }
              )
            }
          >
            View JSON
          </Button>
        }
      >
        {id && (
          <CertPolicyForm policyId={id} value={certPolicy} onChange={mutate} />
        )}
      </Card>
    </>
  );
}
