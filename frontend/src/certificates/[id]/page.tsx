import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Typography } from "antd";
import { useContext, useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { parse } from "uuid";
import { DrawerContext } from "../../admin/contexts/DrawerContext";
import { useNamespace } from "../../admin/contexts/useNamespace";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { NumericDateTime } from "../../components/NumericDateTime";
import {
  CertificateToJSON,
  NamespaceProvider,
  UpdatePendingCertificateRequest,
} from "../../generated/apiv2";
import {
  base64UrlEncodedToStdEncoded,
  toPEMBlock,
} from "../../utils/encodingUtils";
import { useAdminApi } from "../../utils/useCertsApi";

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
  const [searchParams] = useSearchParams();
  const pending = searchParams.get("pending");
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAdminApi();
  const {
    data: cert,
    runAsync,
    mutate,
  } = useRequest(
    async (includeJwk?: boolean) => {
      if (id) {
        return await api?.getCertificate({
          id,
          namespaceId,
          namespaceProvider,
          includeJwk,
          pending: pending === "true" ? true : undefined,
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
          await api?.deleteCertificate({
            id,
            namespaceId,
            namespaceProvider,
          });
        } catch {
          // TODO
        }
      } else {
        const data = await api?.deleteCertificate({
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
    const jwk = cert?.jwk ?? (await runAsync(true))?.jwk;
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
        return api?.addMsEntraKeyCredential({
          id,
          namespaceId,
          namespaceProvider,
        });
      }
    },
    { manual: true }
  );

  const { run: updatePendingCert } = useRequest(
    async (req: UpdatePendingCertificateRequest) => {
      if (!id) {
        return;
      }
      return await api?.updatePendingCertificate({
        id,
        namespaceId,
        namespaceProvider,
        updatePendingCertificateRequest: req,
      });
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
            <dd className="capitalize">{cert?.status}</dd>
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
          {cert?.iat && (
            <div>
              <dt>Issued</dt>
              <dd>
                <NumericDateTime value={cert?.iat} />
              </dd>
            </div>
          )}
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
          {cert?.thumbprint && (
            <div>
              <dt>Thumbprint SHA-1</dt>
              <dd className="font-mono uppercase">{cert?.thumbprint}</dd>
            </div>
          )}
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
      {cert?.pendingAcme && (
        <Card title="ACME">
          <ul>
            {cert.pendingAcme.authorizations?.map((authz) => (
              <li key={authz.url}>
                <ul>
                  {authz.challenges?.map((chan) => (
                    <li key={chan.url}>
                      <div>{chan.type}</div>
                      <div>{chan.dnsRecord}</div>
                      <div>
                        <Button
                          onClick={() => {
                            updatePendingCert({
                              acmeAcceptChallengeUrl: chan.url,
                            });
                          }}
                        >
                          Ready to verify
                        </Button>
                      </div>
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>
          {cert.pendingAcme.authorizations?.every(
            (authz) => authz.status === "valid"
          ) && (
            <Button
              onClick={() => {
                updatePendingCert({
                  acmeOrderCertificate: true,
                });
              }}
              type="primary"
            >
              Order certificate
            </Button>
          )}
        </Card>
      )}
      <Card title="Actions">
        <div className="space-y-4">
          {cert?.status === "issued" &&
            namespaceProvider ===
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
