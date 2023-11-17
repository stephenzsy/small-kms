import { Card, Typography } from "antd";
import { useNamespace } from "../../admin/contexts/NamespaceContextRouteProvider";
import { AdminApi } from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useRequest } from "ahooks";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { CertPolicyForm } from "../../admin/forms/CertPolicyForm";

function useCertificatePolicy(id: string | undefined) {
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    async () => {
      if (id) {
        return api.getCertificatePolicy({
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

export default function CertPolicyPage() {
  const { id } = useParams<{ id: string }>();
  const { data: certPolicy } = useCertificatePolicy(id);
  return (
    <>
      <Typography.Title>Certificate Policy</Typography.Title>
      <Card title="Current Certificate Policy">
        <JsonDataDisplay data={certPolicy} />
      </Card>
      <Card>
        {id && <CertPolicyForm policyId={id} value={certPolicy} />}
      </Card>
    </>
  );
}
