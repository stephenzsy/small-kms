import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Descriptions, Typography } from "antd";
import { DescriptionsItemType } from "antd/es/descriptions";
import { useContext, useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { parse } from "uuid";
import { base64UrlEncodedToStdEncoded, toPEMBlock } from "../utils/encodingUtils";

function formatSerialNumber(certID: string): string {
  const a = parse(certID)
  return a.reduce((accumulator, currentValue) => {
    if (accumulator !== "") {
      accumulator += ":"
    }
    return accumulator + currentValue.toString(16).padStart(2, "0");
  }, "");
}

export default function CertificatePage() {
  const { namespaceKind, namespaceId: namespaceIdentifier } = useContext(NamespaceContext);

  const { certId } = useParams() as { certId: string };

  const adminApi = useAuthedClient(AdminApi);
  const { data: cert } = useRequest(() => {
    return adminApi.getCertificate({
      resourceId: certId,
      namespaceId: namespaceIdentifier,
      namespaceKind: namespaceKind,
    });
  }, {});

  const {
    data: deleted,
    loading: deleteLoading,
    run: deleteCert,
  } = useRequest(
    async () => {
      await adminApi.deleteCertificate({
        resourceId: certId,
        namespaceId: namespaceIdentifier,
        namespaceKind,
      });
      return true;
    },
    { manual: true }
  );

  const [certDownloadBlobInfo, setCertDownloadBlob] = useState<[string, Blob]>()

  const createDownloadLink = useMemoizedFn(() => {
    if (!cert || !cert.jwk.x5c) {
      return;
    }
    const blob = new Blob([cert.jwk.x5c.map((b) => toPEMBlock(base64UrlEncodedToStdEncoded(b), "CERTIFICATE")).join("\n")], {
      type: "application/x-pem-file",
    });
    setCertDownloadBlob([cert.id, blob]);
  })

  const [_certDownloadUrl, setCertDownloadUrl] = useState<string>()
  useEffect(() => {
    if (!certDownloadBlobInfo || !cert) {
      return 
    }
    const [, blob] = certDownloadBlobInfo
    const certDownloadUrl = URL.createObjectURL(blob)
    setCertDownloadUrl(certDownloadUrl)
    return () => {
      URL.revokeObjectURL(certDownloadUrl)
      setCertDownloadBlob(undefined)
    }
  }, [certDownloadBlobInfo])

  const certDownloadUrl = (cert && certDownloadBlobInfo?.[0] === cert.id) ? _certDownloadUrl : undefined

  return (
    <>
      <Typography.Title>Certificate</Typography.Title>
      <Card title="Certificate">
        <dl>
          <div>
            <dt className="font-medium">ID</dt>
            <dd className="font-mono">{cert?.id}</dd>
          </div>
          <div>
            <dt className="font-medium">Serial number</dt>
            <dd className="font-mono">{cert?.id && formatSerialNumber(cert.id)}</dd>
          </div>
          <div>
            <dt className="font-medium">Subject common name (CN)</dt>
            <dd className="font-mono">{cert?.subject.commonName}</dd>
          </div>
          <div>
            <dt className="font-medium">Issued</dt>
            <dd>{cert?.attributes.iat &&
              new Date(cert.attributes.iat * 1000).toString()}</dd>
          </div>
          <div>
            <dt className="font-medium">Expires</dt>
            <dd>{cert?.attributes.exp &&
              new Date(cert.attributes.exp * 1000).toString()}</dd>
          </div>
          <div>
            <dt className="font-medium">Thumbprint SHA-1</dt>
            <dd className="font-mono">{cert?.thumbprint}</dd>
          </div>
          <div>
            <dt className="font-medium">DNS Names</dt>
            <dd>{cert?.subjectAlternativeNames?.dnsNames?.join(", ")}</dd>
          </div>
          <div>
            <dt className="font-medium">IP Addresses</dt>
            <dd>{cert?.subjectAlternativeNames?.ipAddresses?.join(", ")}</dd>
          </div>
        </dl>
      </Card>
      <Card title="Actions">
        {cert && !cert.deleted && !deleted && (<div className="flex items-center gap-4">
          {certDownloadUrl ?
            <Button type="primary" href={certDownloadUrl} download={`${cert.id}.pem`}>Download certificate</Button> :
            <Button type="primary" onClick={createDownloadLink}>Get download link</Button>
          }
          <Button
            danger
            onClick={() => {
              deleteCert();
            }}
          >
            {deleteLoading ? "Deleting...." : "Delete"}
          </Button>
        </div>)}
      </Card>
    </>
  );
}
