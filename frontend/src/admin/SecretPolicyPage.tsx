import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Radio,
  Typography,
} from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import {
  AdminApi,
  ResourceKind,
  SecretGenerateMode,
  SecretPolicyParameters,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { useForm, useWatch } from "antd/es/form/Form";
import { JsonDataDisplay } from "../components/JsonDataDisplay";

type SecretPolicyFormState = {
  identifier: string;
  displayName: string;
  mode: SecretGenerateMode;
  byteLength: number;
};

function SecretPolicyForm({ policyId }: { policyId: string }) {
  const [form] = useForm<SecretPolicyFormState>();
  const newPolicyId = useWatch("identifier", form);
  const generateMode = useWatch("mode", form);
  const { namespaceId: namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { run } = useRequest(
    async (req: SecretPolicyParameters) => {
      await api.putSecretPolicy({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        resourceId: policyId || newPolicyId,
        secretPolicyParameters: req,
      });
    },
    { manual: true }
  );

  const onFinish = useMemoizedFn((values: SecretPolicyFormState) => {
    run({
      mode: values.mode,
      displayName: values.displayName,
      expiryTime: undefined,
      randomLength: values.byteLength,
    });
  });

  return (
    <Form<SecretPolicyFormState>
      form={form}
      layout="vertical"
      initialValues={{
        mode: SecretGenerateMode.SecretGenerateModeServerGeneratedRandom,
        byteLength: 32,
      }}
      onFinish={onFinish}
    >
      {!policyId && (
        <Form.Item<SecretPolicyFormState> label="ID" name="identifier" required>
          <Input placeholder="default" />
        </Form.Item>
      )}
      <Form.Item<SecretPolicyFormState> label="Display name" name="displayName">
        <Input placeholder={policyId || newPolicyId} />
      </Form.Item>
      <Form.Item<SecretPolicyFormState>
        label="Secret generate mode"
        name="mode"
        required
      >
        <Radio.Group
          options={[
            {
              label: "Server generated random",
              value: SecretGenerateMode.SecretGenerateModeServerGeneratedRandom,
            },
            {
              label: "Manual",
              value: SecretGenerateMode.SecretGenerateModeManual,
            },
          ]}
        />
      </Form.Item>
      <Form.Item<SecretPolicyFormState>
        name="byteLength"
        label="Byte length"
        required
      >
        <InputNumber />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

export default function SecretPolicyPage() {
  const { namespaceId: namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const _policyId = useParams().policyId;
  const policyId = _policyId === "_create" ? "" : _policyId;

  const api = useAuthedClient(AdminApi);
  const {
    data,
    run: refresh,
    mutate,
  } = useRequest(
    async () =>
      await api.getSecretPolicy({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        resourceId: policyId!,
      }),
    {
      refreshDeps: [policyId, namespaceIdentifier, namespaceKind],
      ready: !!policyId,
    }
  );

  return (
    <>
      <Typography.Title>
        Certificate Policy: {policyId || "new policy"}
      </Typography.Title>
      <div className="font-mono">
        {namespaceKind}:{namespaceIdentifier}:
        {ResourceKind.ResourceKindCertPolicy}/{policyId}
      </div>
      {/* { <Card title="Certificate list">
        <Table<CertificateRef>
          columns={certListColumns}
          dataSource={issuedCertificates}
          rowKey={(r) => r.id}
        />
      </Card>} */}
      <Card title="Manage secrets">
        {/* <div className="space-y-4">
          <RequestCertificateControl
            certPolicyId={certPolicyId}
            onComplete={refreshCertificates}
          />
          {(namespaceKind === NamespaceKind.NamespaceKindRootCA ||
            namespaceKind === NamespaceKind.NamespaceKindIntermediateCA) &&
            (certPolicyId !== issuerRule?.policyId ? (
              <Button
                onClick={() => {
                  setIssuerRule({ policyId: certPolicyId });
                }}
              >
                Set as issuer policy
              </Button>
            ) : (
              <Tag color="blue">Current issuer policy</Tag>
            ))}
          {namespaceKind === NamespaceKind.NamespaceKindServicePrincipal && (
            <div className="flex gap-4 items-center">
              <Button
                onClick={() => {
                  setEntraClientCred({ policyId: certPolicyId });
                }}
              >
                Set as Microsoft Entra ID client credential policy
              </Button>
              {certPolicyId === entraClientCred?.policyId && (
                <Tag color="blue">
                  Current Microsoft Entra ID client credential policy
                </Tag>
              )}
            </div>
          )}
              </div> */}
      </Card>
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
      {policyId !== undefined && (
        <Card title="Create or update secret policy">
          <SecretPolicyForm policyId={policyId} />
          {/* <CertPolicyForm
          certPolicyId={certPolicyId}
          value={data}
          onChange={onMutate}
        /> */}
        </Card>
      )}
    </>
  );
}
