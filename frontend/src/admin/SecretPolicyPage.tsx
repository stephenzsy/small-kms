import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Radio,
  Table,
  TableColumnType,
  Typography,
} from "antd";
import { useContext, useMemo } from "react";
import { useParams } from "react-router-dom";
import {
  AdminApi,
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
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);

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

function GenerateSecretControl({
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
          generateSecret();
        }}
      >
        Request certificate
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

export default function SecretPolicyPage() {
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
      return await api.getSecretPolicy({
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
        Secret Policy: {policyId || "new policy"}
      </Typography.Title>
      <div className="font-mono">
        {namespaceKind}:{namespaceIdentifier}:
        {ResourceKind.ResourceKindCertPolicy}/{policyId}
      </div>
      <Card title="Secret list">
        <Table<SecretRef>
          columns={secretsTableColumns}
          dataSource={issuedSecrets}
          rowKey="id"
        />
      </Card>
      <Card title="Manage secrets">
        {policyId && (
          <GenerateSecretControl
            policyId={policyId}
            onComplete={refreshSecrets}
          />
        )}
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
