import { useRequest } from "ahooks";
import { Card, Typography } from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import { parse } from "uuid";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";

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
  const { namespaceKind, namespaceId: namespaceIdentifier } =
    useContext(NamespaceContext);

  const { certId } = useParams() as { certId: string };

  const adminApi = useAuthedClient(AdminApi);
  const { data: cert } = useRequest(() => {
    return adminApi.getCertificate({
      resourceId: certId,
      namespaceId: namespaceIdentifier,
      namespaceKind: namespaceKind,
    });
  }, {});

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
            <dd className="font-mono">
              {cert?.id && formatSerialNumber(cert.id)}
            </dd>
          </div>
          <div>
            <dt className="font-medium">Subject common name (CN)</dt>
            <dd className="font-mono">{cert?.subject.commonName}</dd>
          </div>
          <div>
            <dt className="font-medium">Issued</dt>
            <dd>
              {cert?.attributes.iat &&
                new Date(cert.attributes.iat * 1000).toString()}
            </dd>
          </div>
          <div>
            <dt className="font-medium">Expires</dt>
            <dd>
              {cert?.attributes.exp &&
                new Date(cert.attributes.exp * 1000).toString()}
            </dd>
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
    </>
  );
}
