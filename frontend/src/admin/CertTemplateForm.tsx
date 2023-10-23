import { XMarkIcon } from "@heroicons/react/24/outline";
import { useMemoizedFn, useRequest } from "ahooks";
import { Alert, Button, Form, Input, InputNumber } from "antd";
import { useForm } from "antd/es/form/Form";
import React, { useEffect, useMemo } from "react";
import {
  AdminApi,
  CertificateTemplateToJSON,
  CertificateUsage,
  NamespaceKind1 as NamespaceKind,
  PutCertificateTemplateRequest,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { JsonDataDisplay } from "../components/JsonDataDisplay";

type FieldType = {
  subjectCN?: string;
  sanDnsNames?: string[];
  sanEmails?: string[];
  sanIps?: string[];
  validityInMonths?: number;
  keyStorePath?: string;
};

function CertTemplateSanListFormItem({
  label,
  fieldName,
  placeholder,
  addLabel = "+ Add",
}: {
  label: React.ReactNode;
  fieldName: keyof FieldType;
  placeholder?: string;
  addLabel?: string;
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
              {addLabel}
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

  const { data, run, error, loading } = useRequest(
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

  const onReset = useMemoizedFn(() => {
    if (data) {
      form.setFieldsValue({
        subjectCN: data.subjectCommonName,
        sanDnsNames: data.subjectAlternativeNames?.dnsNames,
        sanEmails: data.subjectAlternativeNames?.emails,
        sanIps: data.subjectAlternativeNames?.ipAddresses,
        validityInMonths: data.validityMonths,
        keyStorePath: data.keyStorePath,
      });
    }
  });

  const certUsages = useMemo((): Set<CertificateUsage> | undefined => {
    return namespaceKind == NamespaceKind.NamespaceKindCaRoot
      ? new Set([
          CertificateUsage.CertUsageCA,
          CertificateUsage.CertUsageCARoot,
        ])
      : namespaceKind == NamespaceKind.NamespaceKindCaInt
      ? new Set([CertificateUsage.CertUsageCA])
      : templateId == "default-ms-entra-client-creds"
      ? new Set([
          CertificateUsage.CertUsageServerAuth,
          CertificateUsage.CertUsageClientAuth,
        ])
      : undefined;
  }, []);

  const onFinish = useMemoizedFn((values: FieldType) => {
    run({
      namespaceId,
      namespaceKind,
      templateId,
      certificateTemplateParameters: {
        subjectCommonName: values.subjectCN!,
        subjectAlternativeNames: {
          dnsNames: values.sanDnsNames,
          emails: values.sanEmails,
          ipAddresses: values.sanIps,
        },
        validityMonths: values.validityInMonths,
        usages: [...(certUsages ?? [])],
        keyStorePath: values.keyStorePath,
      },
    });
  });

  const onFinishFailed = (errorInfo: any) => {
    console.log("Failed:", errorInfo);
  };

  useEffect(() => {
    onReset();
  }, [data]);

  return (
    <>
      <JsonDataDisplay
        data={data}
        loading={loading}
        toJson={CertificateTemplateToJSON}
        className="mb-4"
      />
      <Form
        form={form}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        style={{ maxWidth: 600 }}
        initialValues={{}}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        autoComplete="off"
      >
        {error && <Alert message={error.message} type="error" showIcon />}
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
          addLabel="+ Add DNS name"
        />
        <CertTemplateSanListFormItem
          label="IPs"
          fieldName="sanIps"
          placeholder="127.0.0.1"
          addLabel="+ Add IP"
        />
        <CertTemplateSanListFormItem
          label="Emails"
          fieldName="sanEmails"
          placeholder="me@example.com"
          addLabel="+ Add email"
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
        <Form.Item<FieldType> label="Key store path" name="keyStorePath">
          <Input placeholder="example.com" />
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
          <div className="flex gap-2">
            <Button type="primary" htmlType="submit">
              Submit
            </Button>
            <Button htmlType="button" onClick={onReset}>
              Reset
            </Button>
          </div>
        </Form.Item>
      </Form>
    </>
  );
}
