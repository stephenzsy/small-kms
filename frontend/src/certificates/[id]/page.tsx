import { Button, Card, Typography } from "antd";
import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { DrawerContext } from "../../admin/contexts/DrawerContext";
import { useNamespace } from "../../admin/contexts/NamespaceContextRouteProvider";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { AdminApi, CertificateToJSON } from "../../generated/apiv2";
import { useMemoizedFn, useRequest } from "ahooks";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { parse } from "uuid";
import { NumericDateTime } from "../../components/NumericDateTime";
import {
  base64UrlEncodedToStdEncoded,
  toPEMBlock,
} from "../../utils/encodingUtils";

function formatSerialNumber(certID: string): string {
  const a = parse(certID);
  return a.reduce((accumulator, currentValue) => {
    if (accumulator !== "") {
      accumulator += ":";
    }
    return accumulator + currentValue.toString(16).padStart(2, "0");
  }, "");
}

export default function CertificatePage() {
  const { id } = useParams<{ id: string }>();
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const {
    data: cert,
    runAsync,
    mutate,
  } = useRequest(
    async (includeJwk?: boolean) => {
      if (id) {
        return await api.getCertificate({
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

  const { run: deleteCert } = useRequest(
    async () => {
      if (!id) {
        return;
      }
      if (cert?.status == "pending") {
        try {
          await api.deleteCertificate({
            id,
            namespaceId,
            namespaceProvider,
          });
        } catch (e) {}
      } else {
        const data = await api.deleteCertificate({
          id,
          namespaceId,
          namespaceProvider,
        });
        mutate(data);
      }
    },
    { manual: true }
  );

  const [blobUrl, setBlobUrl] = useState<string>();
  useEffect(() => {
    if (blobUrl) {
      return () => {
        URL.revokeObjectURL(blobUrl);
      };
    }
  }, [blobUrl]);

  const getDownloadLink = useMemoizedFn(async () => {
    let jwk = cert?.jwk ?? (await runAsync(true))?.jwk;
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
                openDrawer(
                  <JsonDataDisplay data={cert} toJson={CertificateToJSON} />,
                  {
                    title: "Certificate",
                    size: "large",
                  }
                )
              }
            >
              View JSON
            </Button>
            <Button
              size="small"
              onClick={() => {
                runAsync(true).then((data) =>
                  openDrawer(
                    <JsonDataDisplay data={data} toJson={CertificateToJSON} />,
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
            <dd className="font-mono">{cert?.id}</dd>
          </div>
          <div>
            <dt>Status</dt>
            <dd>{cert?.status}</dd>
          </div>
          <div>
            <dt>Serial Number</dt>
            <dd className="font-mono uppercase">
              {cert?.id && formatSerialNumber(cert.id)}
            </dd>
          </div>
          <div>
            <dt>Subject</dt>
            <dd>{cert?.subject}</dd>
          </div>
          <div>
            <dt>Issued</dt>
            <dd>
              <NumericDateTime value={cert?.iat} />
            </dd>
          </div>
          <div>
            <dt>Not Before</dt>
            <dd>
              <NumericDateTime value={cert?.nbf} />
            </dd>
          </div>
          <div>
            <dt>Expires</dt>
            <dd>
              <NumericDateTime value={cert?.exp} />
            </dd>
          </div>
          <div>
            <dt>Thumbprint SHA-1</dt>
            <dd className="font-mono uppercase">{cert?.thumbprint}</dd>
          </div>
          {cert?.subjectAlternativeNames?.dnsNames && (
            <div>
              <dt>DNS Names</dt>
              <dd>{cert?.subjectAlternativeNames?.dnsNames?.join(", ")}</dd>
            </div>
          )}
          {cert?.subjectAlternativeNames?.ipAddresses && (
            <div>
              <dt>IP Addresses</dt>
              <dd>{cert?.subjectAlternativeNames?.ipAddresses?.join(", ")}</dd>
            </div>
          )}
          {cert?.subjectAlternativeNames?.emails && (
            <div>
              <dt>Emails</dt>
              <dd>{cert?.subjectAlternativeNames?.emails?.join(", ")}</dd>
            </div>
          )}
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
                <a href={blobUrl} download={`${cert?.id}.pem`}>
                  {cert?.id}.pem
                </a>
              </span>
            )}
          </div>
        </div>
      </Card>
      <Card title="Danger zone">
        <div className="space-y-4">
          <Button
            danger
            onClick={() => {
              deleteCert();
            }}
          >
            {cert?.status === "pending"
              ? "Delete certificate"
              : "Deactivate certificate"}
          </Button>
        </div>
      </Card>
    </>
  );
}
