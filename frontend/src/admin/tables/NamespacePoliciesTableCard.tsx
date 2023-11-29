import { Button, Card, Divider, Form, Input, Typography } from "antd";
import { ResourceRefsTable } from "./ResourceRefsTable";
import { useNamespace } from "../contexts/useNamespace";
import { useMemoizedFn, useRequest } from "ahooks";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { AdminApi } from "../../generated/apiv2";
import { Link } from "../../components/Link";
import { useForm } from "antd/es/form/Form";
import { useNavigate } from "react-router-dom";
import { useMemo } from "react";

export type NamespacePoliciesTableCardProps = {
  type: "cert" | "key" | "issuer";
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
      } else if (type === "issuer") {
        return await api.listExternalCertificateIssuers({
          namespaceId,
        });
      }
    },
    {
      refreshDeps: [namespaceProvider, namespaceId, type],
    }
  );

  const linkPrefix = useMemo(() => {
    return `./${
      type === "cert"
        ? "cert-policies"
        : type === "key"
        ? "key-policies"
        : type === "issuer"
        ? "cert-issuers"
        : ""
    }`;
  }, [type]);

  const renderActions = useMemoizedFn((ref) => {
    const linkTo = `./${linkPrefix}/${ref.id}`;
    return <Link to={linkTo}>View</Link>;
  });

  const [formInstance] = useForm<{ policyId: string }>();

  return (
    <Card
      title={
        type === "cert"
          ? "Certificate Policies"
          : type === "key"
          ? "Key Policies"
          : type === "issuer"
          ? "Certificate Issuer"
          : ""
      }
    >
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
          navigate(`./${linkPrefix}/${policyId}`);
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
