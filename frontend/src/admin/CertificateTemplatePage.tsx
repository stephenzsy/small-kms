import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Checkbox,
  Table,
  TableColumnType,
  Typography,
} from "antd";
import React, { useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import { Link } from "../components/Link";
import { WellknownId } from "../constants";
import {
  AdminApi,
  CertificateRef,
  CertificateTemplate,
  CertificateUsage,
  NamespaceKind,
} from "../generated";
import {
  ValueState,
  ValueStateMayBeFixed,
  useFixedValueState,
  useValueState,
} from "../utils/formStateUtils";
import { useAuthedClient } from "../utils/useCertsApi";
import { CertTemplateForm } from "./CertTemplateForm";
import { InputField } from "./InputField";
import { RefTableColumn, RefsTable } from "./RefsTable";

export interface CertificateTemplateFormState {
  issuerNamespaceId: ValueStateMayBeFixed<string>;
  issuerTemplateId: ValueState<string>;
  subjectCN: ValueState<string>;
  validityInMonths: number;
  setValidityInMonths: (value: number) => void;
  keyStorePath: ValueStateMayBeFixed<string>;
  certUsages: ValueStateMayBeFixed<ReadonlySet<CertificateUsage>>;
}

export function useCertificateTemplateFormState(
  certTemplate: CertificateTemplate | undefined,
  nsType: NamespaceKind | undefined,
  nsId: string,
  templateId: string
): CertificateTemplateFormState {
  const [validityInMonths, setValidityInMonths] = React.useState<number>(0);

  const randKeyStoreSuffix = useMemo(() => {
    return uuidv4().substring(0, 8);
  }, [nsType, nsId, templateId]);

  const fixedIssuerNamespaceId = useMemo(() => {
    switch (nsType) {
      case NamespaceKind.NamespaceKindCaRoot:
        return nsId;
      case NamespaceKind.NamespaceKindCaInt:
        return nsId === WellknownId.nsTestIntCa
          ? WellknownId.nsTestRootCa
          : WellknownId.nsRootCa;
    }
    return undefined;
  }, [nsType, nsId]);

  const state: CertificateTemplateFormState = {
    issuerNamespaceId: useFixedValueState(
      useValueState(""),
      fixedIssuerNamespaceId
    ),
    issuerTemplateId: useValueState("default"),
    subjectCN: useValueState(""),
    validityInMonths,
    setValidityInMonths,
    keyStorePath: useFixedValueState(
      useValueState(`${nsType}-${randKeyStoreSuffix}`),
      nsType === NamespaceKind.NamespaceKindGroup ? "" : undefined
    ),
    certUsages: useFixedValueState(
      useValueState(
        (): ReadonlySet<CertificateUsage> =>
          new Set([
            CertificateUsage.CertUsageServerAuth,
            CertificateUsage.CertUsageClientAuth,
          ])
      ),
      nsType == NamespaceKind.NamespaceKindCaRoot
        ? new Set([
            CertificateUsage.CertUsageCA,
            CertificateUsage.CertUsageCARoot,
          ])
        : nsType == NamespaceKind.NamespaceKindCaInt
        ? new Set([CertificateUsage.CertUsageCA])
        : templateId == "default-ms-entra-client-creds"
        ? new Set([
            CertificateUsage.CertUsageServerAuth,
            CertificateUsage.CertUsageClientAuth,
          ])
        : undefined
    ),
  };

  React.useEffect(() => {
    if (certTemplate) {
      // state.issuerNamespaceId.onChange?.(certTemplate.issuer.namespaceId);
      // state.issuerTemplateId.onChange(
      //   certTemplate.issuer.templateId ?? uuidNil
      // );
      state.subjectCN.onChange(certTemplate.subjectCommonName);
      setValidityInMonths(certTemplate.validityMonths ?? 0);
      state.keyStorePath.onChange?.(certTemplate.keyStorePath ?? "");
      state.certUsages.onChange?.(new Set(certTemplate.usages));
    }
  }, [certTemplate]);

  return state;
}

const dateShortFormatter = new Intl.DateTimeFormat("en-US", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

function useColumns(namespaceKind: NamespaceKind, namespaceId: string) {
  return useMemo<TableColumnType<CertificateRef>[]>(
    () => [
      {
        title: "ID",
        render: (r) => <span className="font-mono">{r.id}</span>,
      },
      {
        title: "Fingerprint (SHA-1)",
        render: (r) => <span className="font-mono">{r.thumbprint}</span>,
      },
      {
        title: "Date expires",
        render: (r: CertificateRef) => (
          <time
            className="font-mono tabular-nums"
            dateTime={r.notAfter.toISOString()}
          >
            {dateShortFormatter.format(r.notAfter)}
          </time>
        ),
      },
      {
        title: "Status",
        render: (r: CertificateRef) =>
          r.isIssued ? (r.deleted ? "Disabled" : "Issued") : "Pending",
      },
      {
        title: "Actions",
        render: (r) => (
          <Link
            to={`/admin/${namespaceKind}/${namespaceId}/certificates/${r.id}`}
          >
            View
          </Link>
        ),
      },
    ],
    [namespaceKind, namespaceId]
  );
}

export default function CertificateTemplatePage() {
  const {
    namespaceId,
    templateId,
    profileType: namespaceKind,
  } = useParams() as {
    namespaceId: string;
    templateId: string;
    profileType: NamespaceKind;
  };

  const adminApi = useAuthedClient(AdminApi);

  const { data: issuedCertificates } = useRequest(
    () => {
      return adminApi.listCertificatesByTemplate({
        namespaceId,
        namespaceKind,
        templateId: templateId,
      });
    },
    { refreshDeps: [namespaceId, templateId] }
  );
  const columns = useColumns(namespaceKind, namespaceId);
  return (
    <>
      <Typography.Title>Certificate template</Typography.Title>
      <div className="font-mono">
        {namespaceKind}:{namespaceId}/cert-template:{templateId}
      </div>
      <Card title="Issued certificates">
        <Table<CertificateRef>
          columns={columns}
          dataSource={issuedCertificates}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card>
        <RequestCertificateControl
          namespaceId={namespaceId}
          namespaceKind={namespaceKind}
          templateId={templateId}
        />
      </Card>
      <Card title="Certificate template">
        <CertTemplateForm
          namespaceId={namespaceId}
          namespaceKind={namespaceKind}
          templateId={templateId}
        />
      </Card>
      {namespaceKind === NamespaceKind.NamespaceKindServicePrincipal && (
        <KeyvaultRoleAssignmentsCard
          namespaceId={namespaceId}
          namespaceKind={namespaceKind}
          templateId={templateId}
          adminApi={adminApi}
        />
      )}
    </>
  );
}

function RequestCertificateControl({
  namespaceId,
  namespaceKind,
  templateId,
}: {
  namespaceId: string;
  namespaceKind: NamespaceKind;
  templateId: string;
}) {
  const adminApi = useAuthedClient(AdminApi);
  const [force, setForce] = useState(false);
  const { run: issueCert } = useRequest(
    async (force: boolean) => {
      await adminApi.issueCertificateFromTemplate({
        namespaceId,
        templateId,
        namespaceKind,
        force,
      });
    },
    { manual: true }
  );

  return (
    <div className="flex gap-8 items-center">
      <Button
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

const topRoleDefIds = ["4633458b-17de-408a-b874-0445c86b69e6"];

const roleNames: Record<string, string> = {
  "4633458b-17de-408a-b874-0445c86b69e6": "Key Vault Secrets User",
  "b86a8fe4-44ce-4948-aee5-eccb2c155cd7": "Key Vault Secrets Officer",
};

type KeyVaultRoleAssignmentDisplayItem = {
  id: string;
  defId: string;
  displayName: string;
  assigned: boolean;
  isTop: boolean;
};

const displayNameColumn: RefTableColumn<KeyVaultRoleAssignmentDisplayItem> = {
  columnKey: "displayName",
  header: "Display Name",
  render: (item) => item.displayName,
};

function KeyvaultRoleAssignmentsCard({
  namespaceId,
  namespaceKind,
  templateId,
  adminApi,
}: {
  namespaceId: string;
  namespaceKind: NamespaceKind;
  templateId: string;
  adminApi: AdminApi;
}) {
  const { data: roleAssignments, run: getRoleAssignments } = useRequest(
    () => {
      return adminApi.listKeyVaultRoleAssignments({
        namespaceId,
        namespaceKind,
        templateId: templateId,
      });
    },
    { refreshDeps: [namespaceId, templateId], manual: true }
  );

  const transformedRoleAssignments = useMemo(() => {
    if (!roleAssignments) {
      return undefined;
    }
    const t = roleAssignments.map((r, i): KeyVaultRoleAssignmentDisplayItem => {
      const id = r.name ?? r.id ?? "unknown-" + i;
      const defIdsParts = r.roleDefinitionId?.split("/");
      const defId = defIdsParts ? defIdsParts[defIdsParts.length - 1] : id;
      return {
        id,
        defId,
        displayName: roleNames[defId] ?? defId,
        assigned: true,
        isTop: topRoleDefIds.includes(defId),
      };
    });
    const topAssignments = topRoleDefIds.map(
      (id): KeyVaultRoleAssignmentDisplayItem => {
        const matched = t.find((item) => item.defId === id);
        if (matched) {
          return matched;
        }
        return {
          id: `roleDefId=${id}`,
          defId: id,
          displayName: roleNames[id] || id,
          assigned: false,
          isTop: true,
        };
      }
    );

    return topAssignments.concat(t.filter((item) => !item.isTop));
  }, [roleAssignments]);

  const { run: removeRoleAssignment } = useRequest(
    async (roleAssignmentId: string) => {
      await adminApi.removeKeyVaultRoleAssignment({
        namespaceId,
        namespaceKind,
        templateId,
        roleAssignmentId,
      });
      getRoleAssignments();
    },
    { manual: true }
  );

  const { run: addRoleAssignment } = useRequest(
    async (roleDefId: string) => {
      await adminApi.addKeyVaultRoleAssignment({
        namespaceId,
        namespaceKind,
        templateId,
        roleDefinitionId: roleDefId,
      });
      getRoleAssignments();
    },
    { manual: true }
  );

  return (
    <Card title="Azure role assignments">
      <div>
        <Button type="primary" onClick={getRoleAssignments}>
          Get current assignments
        </Button>
      </div>
      <RefsTable
        items={transformedRoleAssignments}
        title="Role assignments"
        columns={[displayNameColumn]}
        refActions={(ref) => {
          if (ref.assigned) {
            return (
              <Button
                danger
                onClick={() => {
                  removeRoleAssignment(ref.id);
                }}
              >
                Remove Assignment
              </Button>
            );
          } else {
            return (
              <Button
                onClick={() => {
                  addRoleAssignment(ref.defId);
                }}
              >
                Add Assignment
              </Button>
            );
          }
        }}
      />
    </Card>
  );
}
