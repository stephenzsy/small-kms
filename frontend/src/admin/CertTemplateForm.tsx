import { XMarkIcon } from "@heroicons/react/24/outline";
import { Button, Checkbox, Form, Input, InputNumber } from "antd";
import { useForm } from "antd/es/form/Form";
import React from "react";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  AdminApi,
  NamespaceKind,
  PutCertificateTemplateRequest,
} from "../generated";
import { useRequest } from "ahooks";

const onFinish = (values: any) => {
  console.log("Success:", values);
};

const onFinishFailed = (errorInfo: any) => {
  console.log("Failed:", errorInfo);
};

type FieldType = {
  subjectCN?: string;
  sanDnsNames?: string[];
  sanEmails?: string[];
  sanIps?: string[];
  validityInMonths?: number;
  remember?: string;
};

function CertTemplateSanListFormItem({
  label,
  fieldName,
  placeholder,
}: {
  label: React.ReactNode;
  fieldName: keyof FieldType;
  placeholder?: string;
}) {
  return (
    <Form.Item<FieldType> label={label}>
      <Form.List name={fieldName}>
        {(subFields, subOpt) => (
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              rowGap: 16,
            }}
          >
            {subFields.map((subField) => (
              <div className="flex gap-2 items-center" key={subField.key}>
                <Form.Item noStyle name={[subField.name]} className="flex-auto">
                  <Input placeholder={placeholder} />
                </Form.Item>
                <XMarkIcon
                  className="h-em w-em flex-none"
                  onClick={() => {
                    subOpt.remove(subField.name);
                  }}
                />
              </div>
            ))}
            <Button type="dashed" onClick={() => subOpt.add()} block>
              + Add
            </Button>
          </div>
        )}
      </Form.List>
    </Form.Item>
  );
}

export function CertTemplateForm({
  namespaceId,
  namespaceKind,
  templateId,
}: {
  namespaceKind: NamespaceKind;
  namespaceId: string;
  templateId: string;
}) {
  const adminApi = useAuthedClient(AdminApi);
  const [form] = useForm<FieldType>();

  const { data, run } = useRequest(
    async (putReq?: PutCertificateTemplateRequest) => {
      if (putReq) {
        return await adminApi.putCertificateTemplate(putReq);
      } else {
        return await adminApi.getCertificateTemplate({
          namespaceKind,
          namespaceId,
          templateId,
        });
      }
    },
    {}
  );

  return (
    <Form
      form={form}
      labelCol={{ span: 8 }}
      wrapperCol={{ span: 16 }}
      style={{ maxWidth: 600 }}
      initialValues={{ remember: true }}
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      autoComplete="off"
    >
      <Form.Item<FieldType>
        label="Subject common name (CN)"
        name="subjectCN"
        rules={[{ required: true, message: "Please enter a common name!" }]}
      >
        <Input placeholder="example.com" />
      </Form.Item>

      <CertTemplateSanListFormItem
        label="DNS Names"
        fieldName="sanDnsNames"
        placeholder="example.com"
      />
      <CertTemplateSanListFormItem
        label="IPs"
        fieldName="sanIps"
        placeholder="127.0.0.1"
      />
      <CertTemplateSanListFormItem
        label="Emails"
        fieldName="sanEmails"
        placeholder="me@example.com"
      />
      <Form.Item<FieldType> name="validityInMonths" label="Validity">
        <InputNumber
          min={1}
          max={360}
          defaultValue={undefined}
          placeholder="default"
          addonAfter="months"
        />
      </Form.Item>

      <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
        <div className="flex gap-2">
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
        <Button htmlType="button">Reset</Button>
        </div>
      </Form.Item>
    </Form>
  );
}
