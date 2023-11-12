import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Checkbox,
  Form,
  Input,
  InputNumber,
  Radio,
  Table,
  TableColumnType,
  Typography,
} from "antd";
import { useContext, useEffect, useMemo } from "react";
import { useParams } from "react-router-dom";
import {
  AdminApi,
  GenerateJsonWebKeyProperties,
  JsonWebKeyCurveName,
  JsonWebKeyOperation,
  JsonWebKeyType,
  KeyPolicy,
  KeyPolicyParameters,
  ResourceKind,
  SecretGenerateMode,
  SecretPolicyParameters,
  SecretRef,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { useForm, useWatch } from "antd/es/form/Form";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { Link } from "../components/Link";

type KeyPolicyFormState = GenerateJsonWebKeyProperties &
  Omit<KeyPolicyParameters, "keyProperties"> & {
    identifier: string;
  };

function KeyPolicyForm({
  policyId,
  value,
}: {
  policyId: string;
  value?: KeyPolicy;
}) {
  const [form] = useForm<KeyPolicyFormState>();
  const newPolicyId = useWatch("identifier", form);
  const kty = useWatch("kty", form);
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);
  const { run } = useRequest(
    async (req: KeyPolicyParameters) => {
      await api.putKeyPolicy({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        resourceId: policyId || newPolicyId,
        keyPolicyParameters: req,
      });
    },
    { manual: true }
  );

  useEffect(() => {
    if (!value) {
      return;
    }
    form.setFieldsValue({
      displayName: value.displayName,
      kty: value.keyProperties.kty,
      keySize: value.keyProperties.keySize,
      crv: value.keyProperties.crv,
      keyOps: value.keyProperties.keyOps,
      expiryTime: value.expiryTime,
      exportable: value.exportable,
    });
  }, [value]);

  const onFinish = useMemoizedFn((values: KeyPolicyFormState) => {
    console.log(values);
    run({
      displayName: values.displayName,
      keyProperties: {
        kty: values.kty,
        keySize: values.keySize,
        crv: values.crv,
        keyOps: values.keyOps,
      },
      expiryTime: values.expiryTime,
      exportable: values.exportable,
    });
  });

  return (
    <Form<KeyPolicyFormState>
      form={form}
      layout="vertical"
      initialValues={{
        mode: SecretGenerateMode.SecretGenerateModeServerGeneratedRandom,
        byteLength: 32,
      }}
      onFinish={onFinish}
    >
      {!policyId && (
        <Form.Item<KeyPolicyFormState> label="ID" name="identifier" required>
          <Input placeholder="default" />
        </Form.Item>
      )}
      <Form.Item<KeyPolicyFormState> label="Display name" name="displayName">
        <Input placeholder={policyId || newPolicyId} />
      </Form.Item>
      <Form.Item<KeyPolicyFormState> label="Key type" name="kty">
        <Radio.Group>
          <Radio value={JsonWebKeyType.Rsa}>RSA</Radio>
          <Radio value={JsonWebKeyType.Ec}>Elliptic Curve</Radio>
          <Radio value={JsonWebKeyType.Oct}>
            Octet sequence (symmetric keys)
          </Radio>
        </Radio.Group>
      </Form.Item>
      {kty === JsonWebKeyType.Rsa && (
        <Form.Item<KeyPolicyFormState> label="Key size" name="keySize">
          <Radio.Group>
            <Radio value={2048}>2048</Radio>
            <Radio value={3072}>3072</Radio>
            <Radio value={4096}>4096</Radio>
          </Radio.Group>
        </Form.Item>
      )}
      {kty === JsonWebKeyType.Ec && (
        <Form.Item<KeyPolicyFormState> label="Curve name" name="crv">
          <Radio.Group>
            <Radio value={JsonWebKeyCurveName.CurveNameP256}>P-256</Radio>
            <Radio value={JsonWebKeyCurveName.CurveNameP384}>P-384</Radio>
            <Radio value={JsonWebKeyCurveName.CurveNameP521}>P-521</Radio>
          </Radio.Group>
        </Form.Item>
      )}
      <Form.Item<KeyPolicyFormState> label="Key operations" name="keyOps">
        <Radio.Group>
          <Radio value={[JsonWebKeyOperation.Sign, JsonWebKeyOperation.Verify]}>
            Sign/Verify
          </Radio>
          <Radio
            value={[JsonWebKeyOperation.Encrypt, JsonWebKeyOperation.Decrypt]}
          >
            Encrypt/Decrypt
          </Radio>
          <Radio
            value={[JsonWebKeyOperation.WrapKey, JsonWebKeyOperation.UnwrapKey]}
          >
            Wrap/Unwrap key
          </Radio>
        </Radio.Group>
      </Form.Item>
      <Form.Item<KeyPolicyFormState> label="Expiry time" name="expiryTime">
        <Input placeholder="P1M" />
      </Form.Item>
      <Form.Item<KeyPolicyFormState>
        name="exportable"
        valuePropName="checked"
        getValueFromEvent={(e) => e.target.checked}
      >
        <Checkbox>Exportable</Checkbox>
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

function GenerateKeyControl({
  policyId,
  onComplete,
}: {
  policyId: string;
  onComplete?: () => void;
}) {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const { run: generateSecret, loading } = useRequest(
    async () => {
      await adminApi.generateSecret({
        resourceId: policyId,
        namespaceId: namespaceId,
        namespaceKind,
      });
      onComplete?.();
    },
    { manual: true }
  );

  return (
    <div className="flex gap-8 items-center">
      <Button
        loading={loading}
        type="primary"
        onClick={() => {
          //generateSecret();
        }}
      >
        Generate cloud key
      </Button>
    </div>
  );
}

function useSecretsTableColumns(activeIssuerCertificateId: string | undefined) {
  return useMemo<TableColumnType<SecretRef>[]>(
    () => [
      {
        title: "ID",
        render: (r: SecretRef) => (
          <>
            <span className="font-mono">{r.id}</span>
          </>
        ),
      },
      {
        title: "Version",
        render: (r: SecretRef) => {
          return <span className="font-mono">{r.version}</span>;
        },
      },
      // {
      //   title: "Expires",
      //   render: (r: CertificateRef) => {
      //     return (
      //       <span className="font-mono">
      //         {r.attributes.exp &&
      //           dateShortFormatter.format(new Date(r.attributes.exp * 1000))}
      //       </span>
      //     );
      //   },
      // },
      // {
      //   title: "Status",
      //   render: (r: CertificateRef) => {
      //     if (r.deleted) {
      //       return <Tag color="red">Deleted</Tag>;
      //     } else if (!r.thumbprint) {
      //       return <Tag color="yellow">Pending</Tag>;
      //     }
      //     return <Tag color="green">Issued</Tag>;
      //   },
      // },
      {
        title: "Actions",
        render: (r: SecretRef) => <Link to={`../secrets/${r.id}`}>View</Link>,
      },
    ],
    [activeIssuerCertificateId]
  );
}

export default function KeyPolicyPage() {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const _policyId = useParams().policyId;
  const policyId = _policyId === "_create" ? "" : _policyId;

  const api = useAuthedClient(AdminApi);
  const {
    data,
    run: refresh,
    mutate,
  } = useRequest(
    async () => {
      return await api.getKeyPolicy({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        resourceId: policyId!,
      });
    },
    {
      refreshDeps: [policyId, namespaceIdentifier, namespaceKind],
      ready: !!policyId,
    }
  );

  const { data: issuedSecrets, run: refreshSecrets } = useRequest(
    async () => {
      return await api.listSecrets({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        policyId: policyId!,
      });
    },
    {
      refreshDeps: [policyId, namespaceIdentifier, namespaceKind],
      ready: !!policyId,
    }
  );

  const secretsTableColumns = useSecretsTableColumns(undefined);
  return (
    <>
      <Typography.Title>
        Key Policy: {policyId || "new policy"}
      </Typography.Title>
      {policyId && (
        <>
          <div className="font-mono">
            {namespaceKind}:{namespaceIdentifier}:
            {ResourceKind.ResourceKindKeyPolicy}/{policyId}
          </div>
          <Card title="Secret list">
            <Table<SecretRef>
              columns={secretsTableColumns}
              dataSource={issuedSecrets}
              rowKey="id"
            />
          </Card>

          <Card title="Manage key">
            {policyId && (
              <GenerateKeyControl
                policyId={policyId}
                onComplete={refreshSecrets}
              />
            )}
          </Card>
        </>
      )}
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
      {policyId !== undefined && (
        <Card title="Create or update secret policy">
          <KeyPolicyForm policyId={policyId} value={data} />
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
