import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Input,
  TableColumnsType,
  Tag,
  Typography
} from "antd";
import { useContext, useId } from "react";
import { useParams } from "react-router-dom";
import { DrawerContext } from "../../admin/contexts/DrawerContext";
import { useNamespace } from "../../admin/contexts/useNamespace";
import { CertPolicyForm } from "../../admin/forms/CertPolicyForm";
import { CertWebEnroll } from "../../admin/forms/CertWebEnroll";
import { ResourceRefsTable } from "../../admin/tables/ResourceRefsTable";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { Link } from "../../components/Link";
import {
  AdminApi,
  CertificateRef,
  CertificateStatus,
  NamespaceProvider,
} from "../../generated/apiv2";
import { dateShortFormatter } from "../../utils/datetimeUtils";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

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

function useCertificateTableColumns(
  currentIssuerId?: string
): TableColumnsType<CertificateRef> {
  return [
    {
      title: "Status",
      key: "status",
      render: (certRef: CertificateRef) => {
        const { id, status } = certRef;
        return (
          <span className="flex">
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
            {id === currentIssuerId && <Tag color="blue">Issuer</Tag>}
          </span>
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
        if (tsNum && tsNum > 0) {
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
  const { data: currentIssuerResp, run: setIssuer } = useRequest(
    async (certificateIdentifier?: string) => {
      if (!id) {
        return;
      }
      try {
        if (certificateIdentifier) {
          return await api.putCertificatePolicyIssuer({
            id: id,
            namespaceId: namespaceId,
            namespaceProvider: namespaceProvider,
            linkRefFields: {
              linkTo: certificateIdentifier,
            },
          });
        }
        return await api.getCertificatePolicyIssuer({
          id: id,
          namespaceId: namespaceId,
          namespaceProvider: namespaceProvider,
        });
      } catch {
        // TODO
      }
    },
    {
      refreshDeps: [namespaceId, namespaceProvider, id],
      ready:
        namespaceProvider === NamespaceProvider.NamespaceProviderRootCA ||
        namespaceProvider === NamespaceProvider.NamespaceProviderIntermediateCA,
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

  const webEnrollCardId = useId();

  const currentIssuerId = currentIssuerResp?.linkTo?.split("/")[1];
  const certColumns = useCertificateTableColumns(currentIssuerId);
  const viewCert = useMemoizedFn((cert: CertificateRef) => {
    return (
      <div className="flex gap-2 items-center">
        <Link to={`../certificates/${cert.id}`}>View</Link>
        {(namespaceProvider === NamespaceProvider.NamespaceProviderRootCA ||
          namespaceProvider ===
            NamespaceProvider.NamespaceProviderIntermediateCA) && (
          <Button
            type="link"
            onClick={() => {
              setIssuer(`${namespaceProvider}:${namespaceId}:cert/${cert.id}`);
            }}
          >
            Set as issuer
          </Button>
        )}
      </div>
    );
  });
  return (
    <>
      <Typography.Title>Certificate Policy</Typography.Title>
      <div className="font-mono">
        {namespaceProvider}:{namespaceId}:cert-policy/{id}
      </div>
      <Card
        title="Certificates"
        extra={
          <div className="flex items-center gap-4">
            {certPolicy?.allowEnroll && (
              <Typography.Link href={`#${webEnrollCardId}`}>
                Enroll Certificate
              </Typography.Link>
            )}
            {certPolicy?.allowGenerate && (
              <Button
                type="link"
                onClick={generateCertificate}
                loading={generateCertificateLoading}
              >
                Generate Certificate
              </Button>
            )}
          </div>
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
      {certPolicy?.allowEnroll && (
        <Card id={webEnrollCardId} title="Enroll Certificate">
          <CertWebEnroll certPolicy={certPolicy} />
        </Card>
      )}
    </>
  );
}
