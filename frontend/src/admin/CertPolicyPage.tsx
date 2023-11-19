import { XMarkIcon } from "@heroicons/react/24/outline";
import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Form,
  Input,
  Table,
  TableColumnType,
  Tag,
  Typography
} from "antd";
import { useContext, useMemo } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useAppAuthContext } from "../auth/AuthProvider";
import { JsonDataDisplay } from "../components/JsonDataDisplay";
import { Link } from "../components/Link";
import {
  AdminApi,
  CertPolicy,
  CertificateRef,
  JsonWebKeyCurveName,
  JsonWebKeyType,
  NamespaceKind,
  ResourceKind,
  SubjectAlternativeNames
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { CertWebEnroll } from "./CertWebEnroll";
import {
  NamespaceConfigContext,
  NamespaceContext,
} from "./contexts/NamespaceContext";

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
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  const adminApi = useAuthedClient(AdminApi);
  const { data, mutate } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.getCertPolicy({
          namespaceKind,
          resourceId: certPolicyId,
          namespaceId,
        });
      }
      return undefined;
    },
    {
      refreshDeps: [certPolicyId, namespaceId, namespaceKind],
    }
  );

  const { data: issuedCertificates, refresh: refreshCertificates } = useRequest(
    async () => {
      if (certPolicyId) {
        return await adminApi.listCertificates({
          namespaceId,
          namespaceKind,
          policyId: certPolicyId,
        });
      }
      return undefined;
    },
    { refreshDeps: [namespaceId, certPolicyId] }
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
  const { isAdmin } = useAppAuthContext();

  return (
    <>
      <Typography.Title>
        Certificate Policy: {certPolicyId || "new policy"}
      </Typography.Title>
      <div className="font-mono">
        {namespaceKind}:{namespaceId}:{ResourceKind.ResourceKindCertPolicy}/
        {certPolicyId}
      </div>
      <Card title="Certificate list">
        <Table<CertificateRef>
          columns={certListColumns}
          dataSource={issuedCertificates}
          rowKey={(r) => r.id}
        />
      </Card>
      {isAdmin && (
        <Card title="Manage certificates">
          <div className="space-y-4">
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
      )}
      <Card title="Current certificate policy">
        <JsonDataDisplay data={data} />
      </Card>
      <Card title="Certificate web enrollment">
        <CertWebEnroll certPolicy={data} />
      </Card>
    </>
  );
}
