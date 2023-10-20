import { Button, Form, Input } from "antd";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  AdminApi,
  AgentProfile,
  AgentProfileParameters,
  AgentProfileToJSON,
} from "../generated";
import { useRequest } from "ahooks";
import { useForm } from "antd/es/form/Form";
import { useEffect } from "react";
import { JsonDataDisplay } from "../components/JsonDataDisplay";

export function ProvisionAgentForm({ namespaceId }: { namespaceId: string }) {
  const adminApi = useAuthedClient(AdminApi);
  const { data, run, loading } = useRequest(
    async (agentProfileParameters?: AgentProfileParameters) => {
      if (agentProfileParameters) {
        return await adminApi.provisionAgentProfile({
          namespaceId,
          agentProfileParameters,
        });
      }
      return await adminApi.getAgentProfile({
        namespaceId,
      });
    },
    {
      defaultParams: [undefined],
      refreshDeps: [namespaceId],
    }
  );
  const [form] = useForm<AgentProfileParameters>();
  useEffect(() => {
    if (data) {
      form.setFieldsValue({
        msEntraClientCredentialCertificateTemplateId:
          data.msEntraClientCredentialCertificateTemplateId,
      });
    }
  }, [data, form]);
  return (
    <div className="space-y-4">
      <JsonDataDisplay
        data={data}
        toJson={AgentProfileToJSON}
        loading={loading}
      />
      <Form
        form={form}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        style={{ maxWidth: 600 }}
        initialValues={{}}
        onFinish={(values) => {
          if (values.msEntraClientCredentialCertificateTemplateId) {
            run({
              msEntraClientCredentialCertificateTemplateId:
                values.msEntraClientCredentialCertificateTemplateId,
            });
          }
        }}
      >
        <Form.Item<AgentProfileParameters>
          label="Certificate Template ID"
          name={"msEntraClientCredentialCertificateTemplateId"}
          required
        >
          <Input placeholder="Enter a certificate template ID" />
        </Form.Item>
        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
          <Button type="primary" htmlType="submit">
            Provision Agent
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}
