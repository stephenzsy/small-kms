import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Input, TableColumnsType, Tag, Typography } from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import { DrawerContext } from "../../admin/contexts/DrawerContext";
import { useNamespace } from "../../admin/contexts/useNamespace";
import { KeyPolicyForm } from "../../admin/forms/KeyPolicyForm";
import { ResourceRefsTable } from "../../admin/tables/ResourceRefsTable";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { Link } from "../../components/Link";
import {
  AdminApi,
  KeyRef,
  KeyStatus,
  NamespaceProvider,
} from "../../generated/apiv2";
import { dateShortFormatter } from "../../utils/datetimeUtils";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

function useKeyPolicy(id: string | undefined) {
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    async () => {
      if (id) {
        return api.getKeyPolicy({
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

function useKeyTableColumns(
  currentIssuerId?: string
): TableColumnsType<KeyRef> {
  return [
    {
      title: "Status",
      key: "status",
      render: (certRef: KeyRef) => {
        const { id, status } = certRef;
        return (
          <span className="flex">
            <Tag
              className="capitalize"
              color={
                status === KeyStatus.KeyStatusActive
                  ? "green"
                  : status === KeyStatus.KeyStatusInactive
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

export default function KeyPolicyPage() {
  const { id } = useParams<{ id: string }>();
  const { data: keyPolicy, mutate } = useKeyPolicy(id);
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const {
    data: keys,
    refresh: refreshKeys,
    loading: keysLoading,
  } = useRequest(
    async () => {
      return await api.listKeys({
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
  const { run: generateKey, loading: generateKeyLoading } = useRequest(
    async () => {
      if (id && keyPolicy) {
        return await api.generateKey({
          id: id,
          namespaceId: namespaceId,
          namespaceProvider: namespaceProvider,
        });
      }
      refreshKeys();
    },
    { manual: true }
  );

  const currentIssuerId = currentIssuerResp?.linkTo?.split("/")[1];
  const keyColumns = useKeyTableColumns(currentIssuerId);
  const viewKeys = useMemoizedFn((cert: KeyRef) => {
    return (
      <div className="flex gap-2 items-center">
        <Link to={`../keys/${cert.id}`}>View</Link>
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
      <Typography.Title>Key Policy</Typography.Title>
      <div className="font-mono">
        {namespaceProvider}:{namespaceId}:key-policy/{id}
      </div>
      <Card
        title="Keys"
        extra={
          <div className="flex items-center gap-4">
            <Button
              type="link"
              onClick={generateKey}
              loading={generateKeyLoading}
            >
              Generate Key
            </Button>
          </div>
        }
      >
        <ResourceRefsTable<KeyRef>
          loading={keysLoading}
          dataSource={keys}
          extraColumns={keyColumns}
          noDisplayName
          renderActions={viewKeys}
        />
      </Card>
      <Card
        title="Key Policy"
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
                      value={`${namespaceProvider}:${namespaceId}:key-policy/${keyPolicy?.id}`}
                    />
                  </label>
                  <JsonDataDisplay data={keyPolicy} />
                </div>,
                {
                  title: "Key Policy",
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
          <KeyPolicyForm policyId={id} value={keyPolicy} onChange={mutate} />
        )}
      </Card>
      {/* {certPolicy?.allowEnroll && (
        <Card id={webEnrollCardId} title="Enroll Certificate">
          <CertWebEnroll certPolicy={certPolicy} />
        </Card>
      )} */}
    </>
  );
}
