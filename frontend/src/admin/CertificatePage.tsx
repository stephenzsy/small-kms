import { useRequest } from "ahooks";
import { Button, Card, Descriptions, Typography } from "antd";
import { DescriptionsItemType } from "antd/es/descriptions";
import { useContext, useMemo } from "react";
import { useParams } from "react-router-dom";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";

export default function CertificatePage() {
  const { namespaceKind, namespaceIdentifier } = useContext(NamespaceContext);

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

  const certDescItems = useMemo<DescriptionsItemType[] | undefined>(() => {
    if (!cert) {
      return undefined;
    }
    return [
      {
        key: 0,
        label: "ID",
        children: cert.id,
      },
      {
        key: "serialNumber",
        label: "Serial Number",
        children: cert.id,
      },
      {
        key: 1,
        label: "Common name",
        children: cert.subject.commonName,
      },
      {
        key: 2,
        label: "Issued",
        children:
          cert.attributes.iat &&
          new Date(cert.attributes.iat * 1000).toString(),
      },
      {
        key: 3,
        label: "Expires",
        children:
          cert.attributes.exp &&
          new Date(cert.attributes.exp * 1000).toString(),
      },
      {
        key: 4,
        label: "Thumbprint SHA-1 hex",
        children: cert.thumbprint,
      },
      {
        key: 5,
        label: "DNS Names",
        children: cert.subjectAlternativeNames?.dnsNames?.join(", "),
      },
      {
        key: 6,
        label: "IP Addresses",
        children: cert.subjectAlternativeNames?.ipAddresses?.join(", "),
      },
    ];
  }, [cert]);

  return (
    <>
      <Typography.Title>Certificate</Typography.Title>
      <Card title="Certificate">
        <Descriptions items={certDescItems} column={1} />
      </Card>
      <Card title="Actions">
        {cert && !cert.deleted && !deleted && (
          <Button
            danger
            onClick={() => {
              deleteCert();
            }}
          >
            {deleteLoading ? "Deleting...." : "Delete"}
          </Button>
        )}
      </Card>
    </>
  );
}
