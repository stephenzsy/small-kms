import { Card, CardSection, CardTitle } from "../components/Card";
import { Link, useParams } from "react-router-dom";
import { AdminApi, NamespaceKind } from "../generated3";
import { useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi3";
import { useMemo } from "react";

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

  return (
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
  );
}
