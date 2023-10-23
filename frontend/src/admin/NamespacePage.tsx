import { useRequest } from "ahooks";
import { Card, Typography } from "antd";
import React, { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Button } from "../components/Button";
import { Card as CCard, CardSection } from "../components/Card";
import Select, { SelectItem } from "../components/Select";
import {
  AdminApi,
  CertificateTemplateRef,
  LinkedCertificateTemplateUsage,
  NamespaceKind1 as NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { AgentConfigurationForm } from "./AgentConfigurationForm";
import {
  CertificateTemplateTable,
  CertificateTemplatesProvider,
  LinkTemplateForm,
} from "./CertificateTemplates";
import { DeviceServicePrincipalLink } from "./DeviceServicePrincipalLink";
import { InputField } from "./InputField";
import { NamespaceContext } from "./NamespaceContext";
import { RefTableColumn, RefsTable } from "./RefsTable";
import { ProvisionAgentForm } from "./ProvisionAgent";

const subjectCnColumn: RefTableColumn<CertificateTemplateRef> = {
  columnKey: "subjectCommonName",
  header: "Subject Common Name",
  render: (item) => {
    return item.linkTo ? (
      <span>Link to: {item.linkTo}</span>
    ) : (
      item.subjectCommonName
    );
  },
};
const enabledColumn: RefTableColumn<CertificateTemplateRef> = {
  columnKey: "enabled",
  header: "Enabled",
  render: (item) => (!item.deleted && item.updated ? "Yes" : "No"),
};

type UsageSelectItem = SelectItem<LinkedCertificateTemplateUsage>;

const selectItems: UsageSelectItem[] = [
  {
    id: LinkedCertificateTemplateUsage.LinkedCertificateTemplateUsageClientAuthorization,
    name: "Client Authorization",
  },
  {
    id: LinkedCertificateTemplateUsage.LinkedCertificateTemplateUsageMemberDelegatedEnrollment,
    name: "Member Delegated Enrollment",
  },
];

export function CertificateTemplatesList({
  nsType,
  namespaceId,
}: {
  nsType: NamespaceKind;
  namespaceId: string;
}) {
  const adminApi = useAuthedClient(AdminApi);
  const { data, run: refresh } = useRequest(
    async () => {
      return await adminApi.listCertificateTemplates({
        namespaceId,
        namespaceKind: nsType,
      });
    },
    {
      refreshDeps: [nsType, namespaceId],
    }
  );

  const [tempalteLinkTarget, setTemplateLinkTarget] = useState("");
  const [selectedUsage, setSelectedUsage] = useState<UsageSelectItem>(
    selectItems[0]
  );

  const { run: createLink } = useRequest(
    async (
      targetLocator: string,
      selectedUsage: LinkedCertificateTemplateUsage
    ) => {
      targetLocator = targetLocator.trim();
      if (!targetLocator) {
        return;
      }
      await adminApi.createLinkedCertificateTemplate({
        namespaceId,
        namespaceKind: nsType,
        createLinkedCertificateTemplateParameters: {
          targetTemplate: targetLocator,
          usage: selectedUsage,
        },
      });

      refresh();
    },
    { manual: true }
  );

  return (
    <>
      <RefsTable
        items={data}
        title="Certificate Templates"
        columns={[subjectCnColumn, enabledColumn]}
        refActions={(ref) => (
          <Link
            to={`/admin/${nsType}/${namespaceId}/certificate-templates/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
      />
      {nsType === NamespaceKind.NamespaceKindServicePrincipal && (
        <div className="flex gap-4 mt-4 items-end">
          <InputField
            labelContent="Create link to"
            value={tempalteLinkTarget}
            onChange={setTemplateLinkTarget}
          />
          <div>
            <Select
              items={selectItems}
              label="Usage"
              selected={selectedUsage}
              setSelected={setSelectedUsage}
            />
          </div>
          <Button
            onClick={() => {
              createLink(tempalteLinkTarget, selectedUsage.id);
            }}
          >
            Add link
          </Button>
        </div>
      )}
    </>
  );
}

function CommonNamespacePage({ children }: React.PropsWithChildren<{}>) {
  const { nsInfo } = React.useContext(NamespaceContext);
  return (
    <>
      <Typography.Title>
        {nsInfo?.type}: {nsInfo?.displayName}
      </Typography.Title>
      <pre>{nsInfo?.id}</pre>
      {children}
    </>
  );
}

function ApplicationNamespacePage({ namespaceId }: { namespaceId: string }) {
  return (
    <CommonNamespacePage>
      <CertificateTemplatesProvider
        namespaceKind={NamespaceKind.NamespaceKindApplication}
        namespaceId={namespaceId}
      >
        <Card title="Certificate templates">
          <CertificateTemplateTable />
        </Card>
        <Card title="Link template">
          <LinkTemplateForm />
        </Card>
        <Card title="Provision agent">
          <ProvisionAgentForm namespaceId={namespaceId} />
        </Card>
      </CertificateTemplatesProvider>
    </CommonNamespacePage>
  );
}

export default function NamespacePage() {
  const { namespaceId, profileType: nsKind } = useParams() as {
    profileType: NamespaceKind;
    namespaceId: string;
  };

  const { nsInfo } = React.useContext(NamespaceContext);
  if (nsKind === NamespaceKind.NamespaceKindApplication) {
    return <ApplicationNamespacePage namespaceId={namespaceId} />;
  }
  const nsType = nsInfo?.type;
  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{nsType}</div>
      {(nsType === NamespaceKind.NamespaceKindCaRoot ||
        nsType === NamespaceKind.NamespaceKindCaInt ||
        nsType == NamespaceKind.NamespaceKindServicePrincipal ||
        nsType == NamespaceKind.NamespaceKindDevice) && (
        <CCard>
          <CardSection>
            <CertificateTemplatesList
              nsType={nsType}
              namespaceId={namespaceId}
            />
          </CardSection>
        </CCard>
      )}
      {nsType === NamespaceKind.NamespaceKindDevice && (
        <DeviceServicePrincipalLink namespaceId={namespaceId} />
      )}
      {nsType === NamespaceKind.NamespaceKindServicePrincipal && (
        <AgentConfigurationForm
          namespaceId={namespaceId}
          namespaceKind={nsType}
        />
      )}
    </>
  );
}
