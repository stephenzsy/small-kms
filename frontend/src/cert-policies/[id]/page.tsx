import { Button, Card, Form, Input, Typography } from "antd";
import { useNamespace } from "../../admin/contexts/NamespaceContextRouteProvider";
import { AdminApi } from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useRequest } from "ahooks";
import { useParams } from "react-router-dom";
import { JsonDataDisplay } from "../../components/JsonDataDisplay";
import { CertPolicyForm } from "../../admin/forms/CertPolicyForm";
import { useContext } from "react";
import { DrawerContext } from "../../admin/contexts/DrawerContext";

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
  const { data: certPolicy, mutate } = useCertificatePolicy(id);
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const { run: generateCertificate, loading: generateCertificateLoading } =
    useRequest(
      async () => {
        if (id && certPolicy) {
          return await api.generateCertificate({
            id: id,
            namespaceId: namespaceId,
            namespaceProvider: namespaceProvider,
          });
        }
      },
      { manual: true }
    );
  return (
    <>
      <Typography.Title>Certificate Policy</Typography.Title>
      <Card title="Actions">
        <Button
          type="primary"
          onClick={generateCertificate}
          loading={generateCertificateLoading}
        >
          Generate Certificate
        </Button>
      </Card>
      <Card
        title="Certificate Policy"
        extra={
          <Button
            type="link"
            onClick={() =>
              openDrawer(
                <div className="space-y-4">
                  <label>
                    <span className="text-sm mb-2">Policy ID:</span>
                    <Input
                      readOnly
                      className="font-mono"
                      value={`${namespaceProvider}:${namespaceId}:cert-policy/${certPolicy?.id}`}
                    />
                  </label>
                  <JsonDataDisplay data={certPolicy} />
                </div>,
                {
                  title: "Certificate Policy",
                  size: "large",
                }
              )
            }
          >
            View JSON
          </Button>
        }
      >
        {id && (
          <CertPolicyForm policyId={id} value={certPolicy} onChange={mutate} />
        )}
      </Card>
    </>
  );
}
