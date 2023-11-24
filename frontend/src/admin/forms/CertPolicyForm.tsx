import { XMarkIcon } from "@heroicons/react/24/outline";
import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Checkbox,
  Divider,
  Form,
  Input,
  Select,
  Typography,
} from "antd";
import { CheckboxChangeEvent } from "antd/es/checkbox";
import { useForm } from "antd/es/form/Form";
import FormItem from "antd/es/form/FormItem";
import { DefaultOptionType } from "antd/es/select";
import { useEffect, useMemo, useState } from "react";
import {
  AdminApi,
  CertificatePolicy,
  CertificatePolicyParameters,
  NamespaceProvider,
} from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useNamespace } from "../contexts/useNamespace";
import { KeyExportableFormItem, KeySpecFormItems } from "./PolicyFormItems";

function SANFormList({
  name,
  addButtonLabel,
  inputPlaceholder,
}: {
  addButtonLabel: React.ReactNode;
  inputPlaceholder?: string;
  name: string[];
}) {
  return (
    <Form.List name={name}>
      {(subFields, subOpt) => {
        return (
          <div className="flex flex-col gap-4">
            {subFields.map((subField) => (
              <div key={subField.key} className="flex items-center gap-4">
                <Form.Item noStyle name={subField.name} className="flex-auto">
                  <Input placeholder={inputPlaceholder} />
                </Form.Item>
                <Button
                  type="text"
                  onClick={() => {
                    subOpt.remove(subField.name);
                  }}
                >
                  <XMarkIcon className="h-em w-em" />
                </Button>
              </div>
            ))}
            <Button type="dashed" onClick={() => subOpt.add()} block>
              {addButtonLabel}
            </Button>
          </div>
        );
      }}
    </Form.List>
  );
}

function IssuerSelector({
  api,
  namespaceProvider,
  onChange,
  value,
}: {
  api: AdminApi;
  namespaceProvider: NamespaceProvider;
  value?: string | undefined;
  onChange?: (value: string | undefined) => void;
}) {
  const issuerNamespaceProvider =
    namespaceProvider === NamespaceProvider.NamespaceProviderIntermediateCA
      ? NamespaceProvider.NamespaceProviderRootCA
      : NamespaceProvider.NamespaceProviderIntermediateCA;

  const { data } = useRequest(
    () => {
      return api.listProfiles({
        namespaceProvider: issuerNamespaceProvider,
      });
    },
    {
      refreshDeps: [issuerNamespaceProvider],
    }
  );

  const options = useMemo((): DefaultOptionType[] => {
    if (!data) {
      return [];
    }
    return data.map((v) => {
      return {
        label: (
          <span>
            {v.displayName} ({issuerNamespaceProvider}:{v.id})
          </span>
        ),
        value: v.id,
      };
    });
  }, [data, issuerNamespaceProvider]);

  const [selectedNamespaceId, setSelectedNamespaceId] = useState<string>();

  const { data: policies } = useRequest(
    async () => {
      if (selectedNamespaceId) {
        return await api.listCertificatePolicies({
          namespaceId: selectedNamespaceId,
          namespaceProvider: issuerNamespaceProvider,
        });
      }
    },
    {
      refreshDeps: [selectedNamespaceId],
    }
  );

  const policyOptions = useMemo((): DefaultOptionType[] => {
    if (!policies) {
      return [];
    }
    return policies.map((v) => {
      return {
        label: (
          <span>
            {v.displayName} ({v.id})
          </span>
        ),
        value: `${issuerNamespaceProvider}:${selectedNamespaceId}:cert-policy/${v.id}`,
      };
    });
  }, [policies, issuerNamespaceProvider, selectedNamespaceId]);

  useEffect(() => {
    if (value) {
      const nsId = value?.split(":")[1];
      if (nsId) {
        setSelectedNamespaceId(nsId);
      }
    }
  }, [value]);

  return (
    <div>
      <Typography.Title level={4}>Issuer</Typography.Title>

      <Form.Item label="Select issuer namespace">
        <Select
          options={options}
          value={selectedNamespaceId}
          onChange={(value) => setSelectedNamespaceId(value)}
        />
      </Form.Item>
      <Form.Item<CertificatePolicyParameters>
        label="Select issuer policy"
        name="issuerPolicyIdentifier"
      >
        <Select<string> options={policyOptions} onChange={onChange} />
      </Form.Item>
    </div>
  );
}

export function CertPolicyForm({
  policyId,
  value,
  onChange,
}: {
  policyId: string;
  value: CertificatePolicy | undefined;
  onChange?: (value: CertificatePolicy | undefined) => void;
}) {
  const [form] = useForm<CertificatePolicyParameters>();
  const { namespaceId, namespaceProvider } = useNamespace();

  const adminApi = useAuthedClientV2(AdminApi);

  const { run, loading } = useRequest(
    async (id: string, params: CertificatePolicyParameters) => {
      const result = await adminApi.putCertificatePolicy({
        namespaceProvider,
        id,
        namespaceId: namespaceId,
        certificatePolicyParameters: params,
      });
      onChange?.(result);
      return result;
    },
    {
      manual: true,
    }
  );

  const isCA =
    namespaceProvider === NamespaceProvider.NamespaceProviderRootCA ||
    namespaceProvider === NamespaceProvider.NamespaceProviderIntermediateCA;

  const onFinish = useMemoizedFn((values: CertificatePolicyParameters) => {
    run(policyId, values);
  });

  useEffect(() => {
    if (!value) {
      return;
    }
    form.setFieldsValue(value);
  }, [value, form]);

  return (
    <Form<CertificatePolicyParameters>
      form={form}
      layout="vertical"
      initialValues={
        value ?? {
          keyExportable: isCA ? false : true,
          allowGenerate: isCA ? true : false,
          allowEnroll: isCA ? false : true,
        }
      }
      onFinish={onFinish}
    >
      <Form.Item<CertificatePolicyParameters>
        name="displayName"
        label="Display Name"
      >
        <Input />
      </Form.Item>
      {namespaceProvider !== NamespaceProvider.NamespaceProviderRootCA && (
        <>
          <Divider />

          <FormItem<CertificatePolicyParameters>
            noStyle
            name="issuerPolicyIdentifier"
          >
            <IssuerSelector
              api={adminApi}
              namespaceProvider={namespaceProvider}
            />
          </FormItem>
        </>
      )}
      <Divider />
      <div>
        <Typography.Title level={4}>Key Specification</Typography.Title>
        <KeySpecFormItems<CertificatePolicyParameters>
          formInstance={form}
          ktyName={["keySpec", "kty"]}
          keySizeName={["keySpec", "keySize"]}
          crvName={["keySpec", "crv"]}
          keyOpsName={["keySpec", "keyOps"]}
        />
        <KeyExportableFormItem<CertificatePolicyParameters>
          name={"keyExportable"}
        />
      </div>
      <Divider />
      <div>
        <Typography.Title level={4}>Policy Attributes</Typography.Title>
        <Form.Item<CertificatePolicyParameters>
          name={"allowGenerate"}
          valuePropName="checked"
          getValueFromEvent={(e: CheckboxChangeEvent) => {
            return e.target.checked;
          }}
        >
          <Checkbox disabled={isCA}>Allow Generate Certificate</Checkbox>
        </Form.Item>
        <Form.Item<CertificatePolicyParameters>
          name={"allowEnroll"}
          valuePropName="checked"
          getValueFromEvent={(e: CheckboxChangeEvent) => {
            return e.target.checked;
          }}
        >
          <Checkbox disabled={isCA}>Allow Enroll Certificate</Checkbox>
        </Form.Item>
      </div>

      <Divider />
      <div>
        <Typography.Title level={4}>Subject</Typography.Title>
        <Form.Item<CertificatePolicyParameters>
          name={["subject", "cn"]}
          label="Common name (CN)"
          required
        >
          <Input placeholder="example.org" />
        </Form.Item>
      </div>
      <Divider />
      <div>
        <Typography.Title level={4}>Subject Alternative Names</Typography.Title>
        <Form.Item<CertificatePolicyParameters> label="DNS names">
          <SANFormList
            name={["sans", "dnsNames"]}
            addButtonLabel="+ Add DNS name"
            inputPlaceholder="example.com"
          />
        </Form.Item>
        <Form.Item<CertificatePolicyParameters> label="IP addresses">
          <SANFormList
            name={["sans", "ipAddresses"]}
            addButtonLabel="+ Add IP Address"
            inputPlaceholder="127.0.0.1 or ::1"
          />
        </Form.Item>

        <Form.Item<CertificatePolicyParameters> label="Email addresses">
          <SANFormList
            name={["sans", "emails"]}
            addButtonLabel="+ Add Email Address"
            inputPlaceholder="example@example.com"
          />
        </Form.Item>
      </div>

      <Divider />

      <Form.Item<CertificatePolicyParameters>
        name="expiryTime"
        label="Expiry time"
      >
        <Input placeholder="P1Y" />
      </Form.Item>
      <Form.Item>
        <Button htmlType="submit" type="primary" loading={loading}>
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}
