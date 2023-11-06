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
import type { CheckboxChangeEvent } from "antd/es/checkbox";
import { useForm, useWatch } from "antd/es/form/Form";
import { useContext, useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { Link } from "../components/Link";
import {
  AdminApi,
  CertPolicy,
  CertPolicyParameters,
  CertificateRef,
  JsonWebKeyCurveName,
  JsonWebKeyOperation,
  JsonWebKeyType,
  NamespaceKind,
  ProfileRef,
  ResourceKind,
  SubjectAlternativeNames,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { CertificateIssuerNamespaceSelect } from "./CertPolicySelector";
import {
  NamespaceConfigContext,
  NamespaceContext,
} from "./contexts/NamespaceContext";
import { XMarkIcon } from "@heroicons/react/24/outline";

function RequestCertificateControl({
  certPolicyId,
  onComplete,
}: {
  certPolicyId: string;
  onComplete?: () => void;
}) {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const [force, setForce] = useState(false);
  const { run: issueCert, loading } = useRequest(
    async (force: boolean) => {
      await adminApi.createCertificate({
        resourceId: certPolicyId,
        namespaceId: namespaceIdentifier,
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
  keyExportable?: boolean;
  kty: JsonWebKeyType;
  keySize: number;
  crv: JsonWebKeyCurveName;
  sans: SubjectAlternativeNames;
  issuerNamespaceId: string;
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
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);

  const { data: issuerProfiles } = useRequest(
    async (): Promise<ProfileRef[] | null> => {
      if (namespaceKind === NamespaceKind.NamespaceKindRootCA) {
        return null;
      }
      if (namespaceKind === NamespaceKind.NamespaceKindIntermediateCA) {
        return await adminApi.listProfiles({
          profileResourceKind: ResourceKind.ProfileResourceKindRootCA,
        });
      }
      return await adminApi.listProfiles({
        profileResourceKind: ResourceKind.ProfileResourceKindIntermediateCA,
      });
    },
    {
      refreshDeps: [namespaceKind],
    }
  );

  const { run } = useRequest(
    async (name: string, params: CertPolicyParameters) => {
      const result = await adminApi.putCertPolicy({
        namespaceKind: namespaceKind,
        resourceId: name,
        namespaceId: namespaceIdentifier,
        certPolicyParameters: params,
      });
      onChange?.(result);
      return result;
    },
    {
      manual: true,
    }
  );

  const ktyState = useWatch("kty", form);
  const keyExportable = useWatch("keyExportable", form);

  //const _selfSigning = useWatch("selfSigning", form);

  const isSelfSigning = namespaceKind === NamespaceKind.NamespaceKindRootCA;
  // ? true
  // : namespaceKind === NamespaceKind.NamespaceKindIntermediateCA
  // ? false
  // : _selfSigning;
  // */

  const onFinish = useMemoizedFn((values: CertPolicyFormState) => {
    run(certPolicyId || values.certPolicyId, {
      expiryTime: values.expiryTime,
      subject: {
        commonName: values.subjectCN,
      },
      displayName: values.displayName,
      keyExportable: values.keyExportable,
      keySpec: {
        kty: values.kty,
        keySize: values.keySize,
        crv: values.crv,
        keyOps: [
          JsonWebKeyOperation.JsonWebKeyOperationSign,
          JsonWebKeyOperation.JsonWebKeyOperationVerify,
        ],
      },
      subjectAlternativeNames: values.sans,
      issuerNamespaceKind:
        namespaceKind === NamespaceKind.NamespaceKindRootCA ||
        namespaceKind === NamespaceKind.NamespaceKindIntermediateCA
          ? NamespaceKind.NamespaceKindRootCA
          : NamespaceKind.NamespaceKindIntermediateCA,
      issuerNamespaceIdentifier:
        namespaceKind === NamespaceKind.NamespaceKindRootCA
          ? namespaceIdentifier
          : values.issuerNamespaceId,
    });
  });

  useEffect(() => {
    if (!value) {
      return;
    }
    form.setFieldsValue({
      certPolicyId: value.id,
      displayName: value.displayName,
      subjectCN: value.subject.commonName,
      expiryTime: value.expiryTime,
      keyExportable: value.keyExportable,
      kty: value.keySpec.kty,
      keySize: value.keySpec.keySize,
      crv: value.keySpec.crv,
      issuerNamespaceId: value.issuerNamespaceIdentifier,
      sans: value.subjectAlternativeNames ?? {},
    });
  }, [value]);

  return (
    <Form<CertPolicyFormState>
      form={form}
      layout="vertical"
      initialValues={{
        certPolicyId: certPolicyId,
        expiryTime: "",
        subjectCN: "",
        displayName: "",
        kty: JsonWebKeyType.JsonWebKeyTypeRSA,
        keySize: 2048,
        selfSigning: false,
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
      {!isSelfSigning && (
        <Form.Item<CertPolicyFormState>
          name="issuerNamespaceId"
          getValueFromEvent={(v: any) => v}
        >
          <CertificateIssuerNamespaceSelect
            availableNamespaceProfiles={issuerProfiles}
            profileKind={
              namespaceKind === NamespaceKind.NamespaceKindIntermediateCA
                ? ResourceKind.ProfileResourceKindRootCA
                : ResourceKind.ProfileResourceKindIntermediateCA
            }
          />
        </Form.Item>
      )}
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

      <div className="ring-1 ring-neutral-300 p-4 rounded-md space-y-4 mb-6">
        <div className="text-lg font-semibold">Subject alternative names</div>
        <Form.Item<CertPolicyFormState> label="DNS names">
          <SANFormList
            name={["sans", "dnsNames"]}
            addButtonLabel="+ Add DNS name"
            inputPlaceholder="example.com"
          />
        </Form.Item>
        <Form.Item<CertPolicyFormState> label="IP addresses">
          <SANFormList
            name={["sans", "ipAddresses"]}
            addButtonLabel="+ Add IP Address"
            inputPlaceholder="127.0.0.1 or ::1"
          />
        </Form.Item>

        <Form.Item<CertPolicyFormState> label="Email addresses">
          <SANFormList
            name={["sans", "emails"]}
            addButtonLabel="+ Add Email Address"
            inputPlaceholder="example@example.com"
          />
        </Form.Item>
      </div>

      <Form.Item<CertPolicyFormState>
        name="expiryTime"
        label="Expiry time"
        required
      >
        <Input placeholder="P1Y" />
      </Form.Item>

      <div className="flex items-start gap-6">
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
      </div>

      <Form.Item>
        <Button htmlType="submit" type="primary">
          Submit
        </Button>
      </Form.Item>
    </Form>
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

function CertificateActions({
  certRef,
  certPolicyId,
  onSetIssuerPolicy,
}: CertificateActionsProps) {
  return (
    <div className="flex items-center gap-2">
      <Link to={`../cert/${certRef.id}`}>View</Link>
    </div>
  );
}

function useCertTableColumns(activeIssuerCertificateId: string | undefined) {
  return useMemo<TableColumnType<CertificateRef>[]>(
    () => [
      {
        title: "Certificate ID",
        render: (r: CertificateRef) => (
          <>
            <span className="font-mono">{r.id}</span>
            {activeIssuerCertificateId === r.id && (
              <Tag className="ml-2" color="blue">
                Current issuer
              </Tag>
            )}
          </>
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
        title: "Status",
        render: (r: CertificateRef) => {
          if (r.deleted) {
            return <Tag color="red">Deleted</Tag>;
          } else if (!r.thumbprint) {
            return <Tag color="yellow">Pending</Tag>;
          }
          return <Tag color="green">Issued</Tag>;
        },
      },
      {
        title: "Actions",
        render: (r) => <CertificateActions certRef={r} />,
      },
    ],
    [activeIssuerCertificateId]
  );
}

export default function CertPolicyPage() {
  const { certPolicyId: _certPolicyId } = useParams() as {
    certPolicyId: string;
  };
  const navigate = useNavigate();
  const certPolicyId = _certPolicyId === "_create" ? "" : _certPolicyId;
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const { data, mutate } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.getCertPolicy({
          namespaceKind,
          resourceId: certPolicyId,
          namespaceId: namespaceIdentifier,
        });
      }
      return undefined;
    },
    {
      refreshDeps: [certPolicyId, namespaceIdentifier, namespaceKind],
    }
  );

  const { data: issuedCertificates, refresh: refreshCertificates } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.listCertificates({
          namespaceId: namespaceIdentifier,
          namespaceKind,
          policyId: certPolicyId,
        });
      }
      return undefined;
    },
    { refreshDeps: [namespaceIdentifier, certPolicyId] }
  );

  const onMutate = useMemoizedFn((value: CertPolicy | undefined) => {
    mutate(value);
    if (!certPolicyId && value) {
      navigate(`./../${value.id}`, { replace: true });
    } else {
      refreshCertificates();
    }
  });

  const {
    issuer: issuerRule,
    setIssuer: setIssuerRule,
    entraClientCred,
    setEntraClientCred,
  } = useContext(NamespaceConfigContext);

  const certListColumns = useCertTableColumns(issuerRule?.certificateId);
  return (
    <>
      <Typography.Title>
        Certificate Policy: {certPolicyId || "new policy"}
      </Typography.Title>
      <div className="font-mono">
        {namespaceKind}:{namespaceIdentifier}:
        {ResourceKind.ResourceKindCertPolicy}/{certPolicyId}
      </div>
      <Card title="Certificate list">
        <Table<CertificateRef>
          columns={certListColumns}
          dataSource={issuedCertificates}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card title="Manage certificates">
        <div className="space-y-4">
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
        </div>
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
