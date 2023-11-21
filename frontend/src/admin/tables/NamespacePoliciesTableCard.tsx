import { Button, Card, Divider, Form, Input, Typography } from "antd";
import { ResourceRefsTable } from "./ResourceRefsTable";
import { useNamespace } from "../contexts/NamespaceContextRouteProvider";
import { useMemoizedFn, useRequest } from "ahooks";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { AdminApi } from "../../generated/apiv2";
import { Link } from "../../components/Link";
import { useForm } from "antd/es/form/Form";
import { useNavigate } from "react-router-dom";

export type NamespacePoliciesTableCardProps = {
  type: "cert" | "key";
};

export function NamespacePoliciesTableCard({
  type,
}: NamespacePoliciesTableCardProps) {
  const { namespaceProvider, namespaceId } = useNamespace();
  const navigate = useNavigate();
  const api = useAuthedClientV2(AdminApi);
  const { data, loading } = useRequest(
    async () => {
      if (type === "cert") {
        return await api.listCertificatePolicies({
          namespaceId,
          namespaceProvider,
        });
      } else if (type === "key") {
        return await api.listKeyPolicies({
          namespaceId,
          namespaceProvider,
        });
      }
    },
    {
      refreshDeps: [namespaceProvider, namespaceId, type],
    }
  );

  const renderActions = useMemoizedFn((ref) => {
    const linkTo = `./${type}-policies/${ref.id}`;
    return <Link to={linkTo}>View</Link>;
  });

  const [formInstance] = useForm<{ policyId: string }>();

  return (
    <Card title={type === "cert" ? "Certificate Policies" : "Key Policies"}>
      <ResourceRefsTable
        renderActions={renderActions}
        loading={loading}
        dataSource={data}
      />
      <Divider />
      <Typography.Title level={5}>Create policy</Typography.Title>
      <Form
        form={formInstance}
        layout="inline"
        onFinish={(p) => {
          const policyId = p.policyId.trim();
          navigate(`./${type}-policies/${policyId}`);
        }}
      >
        <Form.Item name="policyId" label="Policy ID" required>
          <Input />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Next
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
