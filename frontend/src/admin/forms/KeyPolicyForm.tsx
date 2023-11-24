import { useRequest } from "ahooks";
import { Button, Form, Input } from "antd";
import { useForm } from "antd/es/form/Form";
import { useEffect } from "react";
import {
  AdminApi,
  CreateKeyPolicyRequest,
  JsonWebKeyOperation,
  JsonWebKeyType,
  KeyPolicy,
} from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useNamespace } from "../contexts/useNamespace";
import { KeyExportableFormItem, KeySpecFormItems } from "./PolicyFormItems";

export function KeyPolicyForm({
  policyId,
  value,
  onChange,
}: {
  policyId: string;
  value?: KeyPolicy;
  onChange?: (value: KeyPolicy) => void;
}) {
  const [form] = useForm<CreateKeyPolicyRequest>();
  const { namespaceId, namespaceProvider } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const { run } = useRequest(
    async (req: CreateKeyPolicyRequest) => {
      const result = await api.putKeyPolicy({
        namespaceId,
        namespaceProvider,
        id: policyId,
        createKeyPolicyRequest: req,
      });
      onChange?.(result);
      return result;
    },
    { manual: true }
  );

  useEffect(() => {
    if (!value) {
      return;
    }
    form.setFieldsValue(value);
  }, [value, form]);

  return (
    <Form<CreateKeyPolicyRequest>
      form={form}
      layout="vertical"
      initialValues={{
        keySpec: {
          kty: JsonWebKeyType.Rsa,
          keySize: 2048,
          keyOps: [JsonWebKeyOperation.Sign, JsonWebKeyOperation.Verify],
        },
      }}
      onFinish={run}
    >
      <Form.Item<CreateKeyPolicyRequest>
        label="Display name"
        name="displayName"
      >
        <Input placeholder={policyId} />
      </Form.Item>
      <KeySpecFormItems<CreateKeyPolicyRequest>
        formInstance={form}
        ktyName={["keySpec", "kty"]}
        keySizeName={["keySpec", "keySize"]}
        crvName={["keySpec", "crv"]}
        keyOpsName={["keySpec", "keyOps"]}
      />
      <KeyExportableFormItem<CreateKeyPolicyRequest> name={"exportable"} />
      <Form.Item<CreateKeyPolicyRequest> label="Expiry time" name="expiryTime">
        <Input placeholder="P1M" />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}
