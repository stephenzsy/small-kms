import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Checkbox,
  Form,
  Input,
  Radio,
  Table,
  TableColumnType,
  Tag,
  Typography,
} from "antd";
import { useForm, useWatch } from "antd/es/form/Form";
import { useContext, useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import {
  AdminApi,
  CertPolicy,
  CertPolicyParameters,
  Certificate,
  CertificateRef,
  JsonWebKeyCurveName,
  JsonWebKeyOperation,
  JsonWebKeyType,
  NamespaceKind,
  NamespaceKind1,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./NamespaceContext2";
import { Link } from "../components/Link";

function RequestCertificateControl({ certPolicyId }: { certPolicyId: string }) {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const [force, setForce] = useState(false);
  const { run: issueCert, loading } = useRequest(
    async (force: boolean) => {
      await adminApi.createCertificate({
        namespaceIdentifier: namespaceId,
        namespaceKind,
        resourceIdentifier: certPolicyId,
      });
    },
    { manual: true }
  );

  return (
    <div className="flex gap-8 items-center">
      <Button
        loading={loading}
        type="primary"
        onClick={() => {
          issueCert(force);
        }}
      >
        Request certificate
      </Button>
      <Checkbox
        checked={force}
        onChange={(e) => {
          setForce(e.target.checked);
        }}
      >
        Force
      </Checkbox>
    </div>
  );
}

type CertPolicyFormState = {
  certPolicyId: string;
  displayName: string;
  subjectCN: string;
  expiryTime: string;
  //keyExportable: boolean;
  kty: JsonWebKeyType;
  keySize: number;
  crv: JsonWebKeyCurveName;
};

function CertPolicyForm({
  certPolicyId,
  value,
  onChange,
}: {
  certPolicyId: string;
  value: CertPolicy | undefined;
  onChange?: (value: CertPolicy | undefined) => void;
}) {
  const [form] = useForm<CertPolicyFormState>();
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { run } = useRequest(
    async (name: string, params: CertPolicyParameters) => {
      const result = await adminApi.putCertPolicy({
        namespaceKind: namespaceKind,
        namespaceIdentifier: namespaceId,
        resourceIdentifier: name,
        certPolicyParameters: params,
      });
      onChange?.(result);
      return result;
    },
    {
      manual: true,
    }
  );

  const onFinish = useMemoizedFn((values: CertPolicyFormState) => {
    run(certPolicyId || values.certPolicyId, {
      expiryTime: values.expiryTime,
      subject: {
        commonName: values.subjectCN,
      },
      displayName: values.displayName,
      //      keyExportable: values.keyExportable,
      keySpec: {
        kty: values.kty,
        keySize: values.keySize,
        crv: values.crv,
        keyOps: [
          JsonWebKeyOperation.JsonWebKeyOperationSign,
          JsonWebKeyOperation.JsonWebKeyOperationVerify,
        ],
      },
    });
  });

  useEffect(() => {
    if (!value) {
      return;
    }
    form.setFieldsValue({
      certPolicyId: value.resourceIdentifier,
      displayName: value.displayName,
      subjectCN: value.subject.commonName,
      expiryTime: value.expiryTime,
      // keyExportable: value.keyExportable,
      kty: value.keySpec.kty,
      keySize: value.keySpec.keySize,
      crv: value.keySpec.crv,
    });
  }, [value]);

  const ktyState = useWatch("kty", form);

  return (
    <>
      <Form<CertPolicyFormState>
        form={form}
        layout="vertical"
        initialValues={{
          certPolicyId: certPolicyId,
          expiryTime: "",
          subjectCN: "",
          displayName: "",
          keyExportable: false,
          kty: JsonWebKeyType.JsonWebKeyTypeRSA,
          keySize: 2048,
        }}
        onFinish={onFinish}
      >
        {!certPolicyId && (
          <Form.Item<CertPolicyFormState>
            name="certPolicyId"
            label="Policy ID"
            required
          >
            <Input />
          </Form.Item>
        )}
        <Form.Item<CertPolicyFormState> name="displayName" label="Display name">
          <Input />
        </Form.Item>
        <div className="ring-1 ring-neutral-300 p-4 rounded-md space-y-4 mb-6">
          <div className="text-lg font-semibold">Key specification</div>
          <Form.Item<CertPolicyFormState> name="kty" label="Key type">
            <Radio.Group>
              <Radio value={JsonWebKeyType.JsonWebKeyTypeRSA}>RSA</Radio>
              <Radio value={JsonWebKeyType.JsonWebKeyTypeEC}>EC</Radio>
            </Radio.Group>
          </Form.Item>
          {ktyState === JsonWebKeyType.JsonWebKeyTypeRSA ? (
            <Form.Item<CertPolicyFormState> name="keySize" label="RSA key size">
              <Radio.Group>
                <Radio value={2048}>2048</Radio>
                <Radio value={3072}>3072</Radio>
                <Radio value={4096}>4096</Radio>
              </Radio.Group>
            </Form.Item>
          ) : ktyState === JsonWebKeyType.JsonWebKeyTypeEC ? (
            <Form.Item<CertPolicyFormState> name="crv" label="EC curve name">
              <Radio.Group>
                <Radio value={JsonWebKeyCurveName.JsonWebKeyCurveNameP256}>
                  P-256
                </Radio>
                <Radio value={JsonWebKeyCurveName.JsonWebKeyCurveNameP384}>
                  P-384
                </Radio>
                <Radio value={JsonWebKeyCurveName.JsonWebKeyCurveNameP521}>
                  P-521
                </Radio>
              </Radio.Group>
            </Form.Item>
          ) : null}
        </div>
        <div className="ring-1 ring-neutral-300 p-4 rounded-md space-y-4 mb-6">
          <div className="text-lg font-semibold">Subject</div>
          <Form.Item<CertPolicyFormState>
            name="subjectCN"
            label="Common name (CN)"
            required
          >
            <Input placeholder="example.org" />
          </Form.Item>
        </div>
        <Form.Item<CertPolicyFormState>
          name="expiryTime"
          label="Expiry time"
          required
        >
          <Input placeholder="P1Y" />
        </Form.Item>

        {/*
        <Form.Item<CertPolicyFormState>
          name="keyExportable"
          valuePropName="checked"
        >
          <Checkbox>Key exportable</Checkbox>
        </Form.Item> */}

        <Form.Item>
          <Button htmlType="submit" type="primary">
            Submit
          </Button>
        </Form.Item>
      </Form>
    </>
  );
}

const dateShortFormatter = new Intl.DateTimeFormat("en-US", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

type CertificateActionsProps = {
  certRef: CertificateRef;
  onSetIssuerPolicy?: (policyId: string) => void;
  certPolicyId?: string;
};

function CertificateActions({
  certRef,
  certPolicyId,
  onSetIssuerPolicy,
}: CertificateActionsProps) {
  return (
    <div className="flex items-center gap-2">
      <Link to={`../cert/${certRef.resourceIdentifier}`}>View</Link>
      {certPolicyId && certRef.issuerForPolicy ? (
        <Tag color="success">Current Issuer</Tag>
      ) : (
        onSetIssuerPolicy && (
          <Button
            type="link"
            onClick={() => {
              onSetIssuerPolicy(certRef.resourceIdentifier);
            }}
          >
            Set as issuer
          </Button>
        )
      )}
    </div>
  );
}

function useColumns(restActionProps: Omit<CertificateActionsProps, "certRef">) {
  return useMemo<TableColumnType<CertificateRef>[]>(
    () => [
      {
        title: "Certificate ID",
        render: (r: CertificateRef) => (
          <span className="font-mono">{r.resourceIdentifier}</span>
        ),
      },
      {
        title: "Thumbprint (SHA-1)",
        render: (r: CertificateRef) => {
          return <span className="font-mono">{r.thumbprint}</span>;
        },
      },
      {
        title: "Expires",
        render: (r: CertificateRef) => {
          return (
            <span className="font-mono">
              {r.attributes.exp &&
                dateShortFormatter.format(new Date(r.attributes.exp * 1000))}
            </span>
          );
        },
      },
      {
        title: "Actions",
        render: (r) => <CertificateActions certRef={r} {...restActionProps} />,
      },
    ],
    [restActionProps]
  );
}

export default function CertPolicyPage() {
  const { certPolicyId: _certPolicyId } = useParams() as {
    certPolicyId: string;
  };
  const navigate = useNavigate();
  const certPolicyId = _certPolicyId === "_create" ? "" : _certPolicyId;
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const { data, mutate } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.getCertPolicy({
          namespaceKind,
          namespaceIdentifier: namespaceId,
          resourceIdentifier: certPolicyId,
        });
      }
      return undefined;
    },
    {
      refreshDeps: [certPolicyId, namespaceId, namespaceKind],
    }
  );
  const onMutate = useMemoizedFn((value: CertPolicy | undefined) => {
    mutate(value);
    if (!certPolicyId && value) {
      navigate(`./../${value.resourceIdentifier}`, { replace: true });
    }
  });

  const { data: issuedCertificates, refresh: refreshCertificate } = useRequest(
    () => {
      return adminApi.listCertificates({
        namespaceIdentifier: namespaceId,
        namespaceKind,
        policyId: certPolicyId,
      });
    },
    { refreshDeps: [namespaceId, certPolicyId] }
  );

  const { run: setIssuerPolicy } = useRequest(
    async (issuerId: string) => {
      await adminApi.setIssuerCertificate({
        namespaceIdentifier: namespaceId,
        namespaceKind: namespaceKind,
        resourceIdentifier: certPolicyId,
        policyIssuerCertRequest: {
          issuerId,
        },
      });
      refreshCertificate();
    },
    {
      manual: true,
    }
  );

  const restProps = useMemo(() => {
    switch (namespaceKind) {
      case NamespaceKind.NamespaceKindRootCA:
      case NamespaceKind.NamespaceKindIntermediateCA:
        return {
          onSetIssuerPolicy: setIssuerPolicy,
          certPolicyId,
        };
    }
    return {};
  }, [namespaceKind]);

  const columns = useColumns(restProps);
  return (
    <>
      <Typography.Title>
        Certificate Policy: {certPolicyId || "new policy"}
      </Typography.Title>
      <div>
        {namespaceKind}/{namespaceId}
      </div>
      <Card title="Certificate list">
        <Table<CertificateRef>
          columns={columns}
          dataSource={issuedCertificates}
          rowKey={(r) => r.resourceIdentifier}
        />
      </Card>
      <Card title="Manage certificates">
        <RequestCertificateControl certPolicyId={certPolicyId} />
      </Card>
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
      <Card title="Create or update certificate policy">
        <CertPolicyForm
          certPolicyId={certPolicyId}
          value={data}
          onChange={onMutate}
        />
      </Card>
    </>
  );
}
