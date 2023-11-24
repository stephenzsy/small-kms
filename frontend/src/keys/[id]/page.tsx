import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Typography } from "antd";
import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { DrawerContext } from "../../admin/contexts/DrawerContext";
import { useNamespace } from "../../admin/contexts/useNamespace";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { NumericDateTime } from "../../components/NumericDateTime";
import { AdminApi, KeyToJSON, NamespaceProvider } from "../../generated/apiv2";
import {
  base64UrlEncodedToStdEncoded,
  toPEMBlock,
} from "../../utils/encodingUtils";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

export default function KeyPage() {
  const { id } = useParams<{ id: string }>();
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const {
    data: key,
    runAsync,
    //  mutate,
  } = useRequest(
    async (includeJwk?: boolean) => {
      if (id) {
        return await api.getKey({
          id,
          namespaceId,
          namespaceProvider,
          includeJwk,
        });
      }
    },
    {
      refreshDeps: [id, namespaceId, namespaceProvider],
    }
  );

  // const { run: deleteCert } = useRequest(
  //   async () => {
  //     if (!id) {
  //       return;
  //     }
  //     if (cert?.status == "pending") {
  //       try {
  //         await api.deleteCertificate({
  //           id,
  //           namespaceId,
  //           namespaceProvider,
  //         });
  //       } catch {
  //         // TODO
  //       }
  //     } else {
  //       const data = await api.deleteCertificate({
  //         id,
  //         namespaceId,
  //         namespaceProvider,
  //       });
  //       mutate(data);
  //     }
  //   },
  //   { manual: true }
  // );

  const [blobUrl, setBlobUrl] = useState<string>();
  useEffect(() => {
    if (blobUrl) {
      return () => {
        URL.revokeObjectURL(blobUrl);
      };
    }
  }, [blobUrl]);

  const getDownloadLink = useMemoizedFn(async () => {
    const jwk = key?.jwk ?? (await runAsync(true))?.jwk;
    if (!jwk) {
      return;
    }
    const pemString = jwk.x5c
      ?.map((b) => toPEMBlock(base64UrlEncodedToStdEncoded(b), "CERTIFICATE"))
      .join("\n");
    if (!pemString) {
      return;
    }
    const blob = new Blob([pemString], {
      type: "application/x-pem-file",
    });
    setBlobUrl(URL.createObjectURL(blob));
  });

  const { run: installAsMsEntraCredential } = useRequest(
    async () => {
      if (id) {
        return api.addMsEntraKeyCredential({
          id,
          namespaceId,
          namespaceProvider,
        });
      }
    },
    { manual: true }
  );
  return (
    <>
      <Typography.Title>Certificate</Typography.Title>
      <div className="font-mono">
        {namespaceProvider}:{namespaceId}:cert/{id}
      </div>
      <Card
        title="Certificate Information"
        extra={
          <span className="inline-flex gap-4 items-center">
            <Button
              size="small"
              onClick={() =>
                openDrawer(<JsonDataDisplay data={key} toJson={KeyToJSON} />, {
                  title: "Certificate",
                  size: "large",
                })
              }
            >
              View JSON
            </Button>
            <Button
              size="small"
              onClick={() => {
                runAsync(true).then((data) =>
                  openDrawer(
                    <JsonDataDisplay data={data} toJson={KeyToJSON} />,
                    {
                      title: "Certificate",
                      size: "large",
                    }
                  )
                );
              }}
            >
              View Full JSON
            </Button>
          </span>
        }
      >
        <dl className="dl">
          <div>
            <dt>ID</dt>
            <dd className="font-mono">{key?.id}</dd>
          </div>
          <div>
            <dt>Status</dt>
            <dd className="capitalize">{key?.status}</dd>
          </div>
          {key?.iat && (
            <div>
              <dt>Created</dt>
              <dd>
                <NumericDateTime value={key?.iat} />
              </dd>
            </div>
          )}
          <div>
            <dt>Not Before</dt>
            <dd>
            <dd>{key?.nbf ? <NumericDateTime value={key?.nbf} /> : "Never"}</dd>
            </dd>
          </div>
          <div>
            <dt>Expires</dt>
            <dd>{key?.exp ? <NumericDateTime value={key?.exp} /> : "Never"}</dd>
          </div>
        </dl>
      </Card>
      <Card title="Download certificate">
        <div className="space-y-4">
          <div>
            <Button type="primary" onClick={getDownloadLink}>
              Get Download Link (.pem)
            </Button>
          </div>
          <div>
            {blobUrl && (
              <span>
                Download link:{" "}
                <a href={blobUrl} download={`${key?.id}.pem`}>
                  {key?.id}.pem
                </a>
              </span>
            )}
          </div>
        </div>
      </Card>
      <Card title="Actions">
        <div className="space-y-4">
          {namespaceProvider ===
            NamespaceProvider.NamespaceProviderServicePrincipal && (
            <Button
              onClick={() => {
                installAsMsEntraCredential();
              }}
            >
              Install as Microsoft Entra Key Credential
            </Button>
          )}
        </div>
      </Card>
      <Card title="Danger zone">
        <div className="space-y-4">
          {/* <Button
            danger
            onClick={() => {
              deleteCert();
            }}
          >
            {key?.status === "pending"
              ? "Delete certificate"
              : "Deactivate certificate"}
          </Button> */}
        </div>
      </Card>
    </>
  );
}
