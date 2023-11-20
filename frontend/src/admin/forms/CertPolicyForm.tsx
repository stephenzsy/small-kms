import { useBoolean, useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Checkbox,
  Divider,
  Form,
  Input,
  Radio,
  Select,
  Typography,
} from "antd";
import { FormInstance, useForm, useWatch } from "antd/es/form/Form";
import { useEffect, useMemo, useState } from "react";
import { JsonWebKeyCurveName, JsonWebKeyType } from "../../generated";
import {
  AdminApi,
  CertificatePolicy,
  CertificateSubject,
  CreateCertificatePolicyRequest,
  JsonWebKeyOperation,
  NamespaceProvider,
} from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useNamespace } from "../contexts/NamespaceContextRouteProvider";
import { XMarkIcon } from "@heroicons/react/24/outline";
import { CheckboxChangeEvent } from "antd/es/checkbox";
import { DefaultOptionType } from "antd/es/select";
import FormItem from "antd/es/form/FormItem";

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
  }, [data]);

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

  const [selfSigned, setSelfSigned] = useState(false);

  useEffect(() => {
    if (value === "self") {
      setSelfSigned(true);
    } else if (value) {
      const nsId = value?.split(":")[1];
      if (nsId) {
        setSelectedNamespaceId(nsId);
      }
    }
  }, [value]);

  return (
    <div>
      <Typography.Title level={4}>Issuer</Typography.Title>
      <Form.Item>
        <Checkbox
          checked={selfSigned}
          onChange={(e) => {
            if (e.target.checked) {
              onChange?.("self");
            } else {
              onChange?.(undefined);
            }
            setSelfSigned(e.target.checked);
          }}
        >
          Self-Signed
        </Checkbox>
      </Form.Item>
      <Form.Item label="Select issuer namespace">
        <Select
          disabled={selfSigned}
          options={options}
          value={selectedNamespaceId}
          onChange={(value) => setSelectedNamespaceId(value)}
        />
      </Form.Item>
      <Form.Item<CreateCertificatePolicyRequest>
        label="Select issuer policy"
        name="issuerPolicyIdentifier"
      >
        <Select<string>
          options={policyOptions}
          onChange={onChange}
          disabled={selfSigned}
        />
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
  const [form] = useForm<CreateCertificatePolicyRequest>();
  const { namespaceId, namespaceProvider } = useNamespace();

  const adminApi = useAuthedClientV2(AdminApi);

  // const { data: issuerProfiles } = useRequest(
  //   async (): Promise<ProfileRef[] | null> => {
  //     if (namespaceProvider === NamespaceProvider.NamespaceKindRootCA) {
  //       return null;
  //     }
  //     if (namespaceKind === NamespaceKind.NamespaceKindIntermediateCA) {
  //       return await adminApi.listProfiles({
  //         profileResourceKind: ResourceKind.ProfileResourceKindRootCA,
  //       });
  //     }
  //     return await adminApi.listProfiles({
  //       profileResourceKind: ResourceKind.ProfileResourceKindIntermediateCA,
  //     });
  //   },
  //   {
  //     refreshDeps: [namespaceKind],
  //   }
  // );

  const { run, loading } = useRequest(
    async (id: string, params: CreateCertificatePolicyRequest) => {
      const result = await adminApi.putCertificatePolicy({
        namespaceProvider,
        id,
        namespaceId: namespaceId,
        createCertificatePolicyRequest: params,
      });
      onChange?.(result);
      return result;
    },
    {
      manual: true,
    }
  );

  const ktyState = useWatch(["keySpec", "kty"], form);

  const isCA =
    namespaceProvider === NamespaceProvider.NamespaceProviderRootCA ||
    namespaceProvider === NamespaceProvider.NamespaceProviderIntermediateCA;
  //const _selfSigning = useWatch("selfSigning", form);
  //const isSelfSigning = namespaceKind === NamespaceKind.NamespaceKindRootCA;
  // ? true
  // : namespaceKind === NamespaceKind.NamespaceKindIntermediateCA
  // ? false
  // : _selfSigning;
  // */
  const onFinish = useMemoizedFn((values: CreateCertificatePolicyRequest) => {
    run(policyId, values);
  });

  useEffect(() => {
    if (!value) {
      return;
    }
    form.setFieldsValue(value);
  }, [value]);

  return (
    <Form<CreateCertificatePolicyRequest>
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
      <Form.Item<CreateCertificatePolicyRequest>
        name="displayName"
        label="Display Name"
      >
        <Input />
      </Form.Item>
      {namespaceProvider !== NamespaceProvider.NamespaceProviderRootCA && (
        <>
          <Divider />

          <FormItem<CreateCertificatePolicyRequest>
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
        <Form.Item<CreateCertificatePolicyRequest>
          name={["keySpec", "kty"]}
          label="Key Type"
        >
          <Radio.Group>
            <Radio value={JsonWebKeyType.Rsa}>RSA</Radio>
            <Radio value={JsonWebKeyType.Ec}>EC</Radio>
          </Radio.Group>
        </Form.Item>
        {ktyState === JsonWebKeyType.Rsa ? (
          <Form.Item<CreateCertificatePolicyRequest>
            name={["keySpec", "keySize"]}
            label="Key Size"
          >
            <Radio.Group className="flex flex-col gap-4">
              <Radio value={2048}>2048</Radio>
              <Radio value={3072}>3072</Radio>
              <Radio value={4096}>4096</Radio>
            </Radio.Group>
          </Form.Item>
        ) : ktyState === JsonWebKeyType.Ec ? (
          <Form.Item<CreateCertificatePolicyRequest>
            name={["keySpec", "crv"]}
            label="Elliptic Curve Name"
          >
            <Radio.Group className="flex flex-col gap-4">
              <Radio value={JsonWebKeyCurveName.CurveNameP256}>P-256</Radio>
              <Radio value={JsonWebKeyCurveName.CurveNameP384}>P-384</Radio>
              <Radio value={JsonWebKeyCurveName.CurveNameP521}>P-521</Radio>
            </Radio.Group>
          </Form.Item>
        ) : null}
        <Form.Item<CreateCertificatePolicyRequest>
          name={["keySpec", "keyOps"]}
          label="Key Operations"
        >
          <Checkbox.Group className="flex flex-col gap-4">
            <Checkbox value={JsonWebKeyOperation.Sign}>Sign</Checkbox>
            <Checkbox value={JsonWebKeyOperation.Verify}>Verify</Checkbox>
            <Checkbox value={JsonWebKeyOperation.Encrypt}>Encrypt</Checkbox>
            <Checkbox value={JsonWebKeyOperation.Decrypt}>Decrypt</Checkbox>
            <Checkbox value={JsonWebKeyOperation.WrapKey}>Wrap Key</Checkbox>
            <Checkbox value={JsonWebKeyOperation.UnwrapKey}>
              Unwrap Key
            </Checkbox>
            <Checkbox value={JsonWebKeyOperation.DeriveKey}>
              Derive Key
            </Checkbox>
            <Checkbox value={JsonWebKeyOperation.DeriveBits}>
              Derive Bits
            </Checkbox>
          </Checkbox.Group>
        </Form.Item>
        <Form.Item<CreateCertificatePolicyRequest>
          name={"keyExportable"}
          valuePropName="checked"
          getValueFromEvent={(e: CheckboxChangeEvent) => {
            return e.target.checked;
          }}
        >
          <Checkbox disabled={isCA}>Private Key Exportable</Checkbox>
        </Form.Item>
      </div>
      <Divider />
      <div>
        <Typography.Title level={4}>Policy Attributes</Typography.Title>
        <Form.Item<CreateCertificatePolicyRequest>
          name={"allowGenerate"}
          valuePropName="checked"
          getValueFromEvent={(e: CheckboxChangeEvent) => {
            return e.target.checked;
          }}
        >
          <Checkbox disabled={isCA}>Allow Generate Certificate</Checkbox>
        </Form.Item>
        <Form.Item<CreateCertificatePolicyRequest>
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
        <Form.Item<CreateCertificatePolicyRequest>
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
        <Form.Item<CreateCertificatePolicyRequest> label="DNS names">
          <SANFormList
            name={["sans", "dnsNames"]}
            addButtonLabel="+ Add DNS name"
            inputPlaceholder="example.com"
          />
        </Form.Item>
        <Form.Item<CreateCertificatePolicyRequest> label="IP addresses">
          <SANFormList
            name={["sans", "ipAddresses"]}
            addButtonLabel="+ Add IP Address"
            inputPlaceholder="127.0.0.1 or ::1"
          />
        </Form.Item>

        <Form.Item<CreateCertificatePolicyRequest> label="Email addresses">
          <SANFormList
            name={["sans", "emails"]}
            addButtonLabel="+ Add Email Address"
            inputPlaceholder="example@example.com"
          />
        </Form.Item>
      </div>

      <Divider />

      <Form.Item<CreateCertificatePolicyRequest>
        name="expiryTime"
        label="Expiry time"
      >
        <Input placeholder="P1Y" />
      </Form.Item>

      {/* <div className="flex items-start gap-6">
        <Form.Item<CertPolicyFormState>
          name="keyExportable"
          valuePropName="checked"
          getValueFromEvent={(e: CheckboxChangeEvent) => {
            if (e.target.indeterminate) {
              return undefined;
            }
            return e.target.checked;
          }}
        >
          <Checkbox indeterminate={keyExportable === undefined}>
            Key exportable:{" "}
            {keyExportable === undefined ? "default" : keyExportable.toString()}
          </Checkbox>
        </Form.Item>
        {keyExportable !== undefined && (
          <Button
            type="link"
            onClick={() => {
              form.setFieldValue("keyExportable", undefined);
            }}
          >
            Reset to default
          </Button>
        )}
      </div> */}

      <Form.Item>
        <Button htmlType="submit" type="primary" loading={loading}>
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}
