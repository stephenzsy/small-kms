import { Link, useParams } from "react-router-dom";
import { AdminApi, NamespaceKind } from "../generated";
import { useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi";
import { useMemo } from "react";
import { Button, Card, Descriptions, DescriptionsProps } from "antd";
import { DescriptionsItemType } from "antd/es/descriptions";

export default function CertificatePage() {
  const {
    profileType: namespaceKind,
    namespaceId,
    certId,
  } = useParams() as {
    profileType: NamespaceKind;
    namespaceId: string;
    certId: string;
  };

  const adminApi = useAuthedClient(AdminApi);
  const { data: cert } = useRequest(() => {
    return adminApi.getCertificate({
      certificateId: certId,
      namespaceId,
      namespaceKind,
    });
  }, {});

  const {
    data: deleted,
    loading: deleteLoading,
    run: deleteCert,
  } = useRequest(
    async () => {
      await adminApi.deleteCertificate({
        certificateId: certId,
        namespaceId,
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
        key: 7,
        label: "Status",
        children: cert.deleted
          ? "Deleted"
          : cert.isIssued
          ? "Issued"
          : "Pending",
      },
      {
        key: 1,
        label: "Common name",
        children: cert.subjectCommonName,
      },
      {
        key: 2,
        label: "Issued",
        children: cert.notBefore.toString(),
      },
      {
        key: 3,
        label: "Expires",
        children: cert.notAfter.toString(),
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
