import { Card, CardSection, CardTitle } from "../components/Card";
import { Link, useParams } from "react-router-dom";
import { AdminApi, NamespaceKind } from "../generated3";
import { useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi3";
import { useMemo } from "react";
import { Button } from "../components/Button";

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

  return (
    <>
      <Card>
        <CardTitle description={cert?.id}>Certificate</CardTitle>
        <dl>
          <CardSection>
            <dt className="font-medium">Common name</dt>
            <dd>{cert?.subjectCommonName}</dd>
          </CardSection>
          <CardSection>
            <dt className="font-medium">Expires</dt>
            <dd>{cert?.notAfter.toString()}</dd>
          </CardSection>
          <CardSection>
            <dt className="font-medium">Thumbprint SHA-1 hex</dt>
            <dd>{cert?.thumbprint}</dd>
          </CardSection>
        </dl>
      </Card>
      <Card>
        <CardTitle>Actions</CardTitle>
        <CardSection>
          {cert && !cert.isIssued && !deleted && (
            <Button
              variant="soft"
              color="danger"
              onClick={() => {
                deleteCert();
              }}
            >
              {deleteLoading ? "Deleting...." : "Delete"}
            </Button>
          )}
        </CardSection>
      </Card>
    </>
  );
}
