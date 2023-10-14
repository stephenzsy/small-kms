import { useRequest } from "ahooks";
import React, { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Card, CardSection } from "../components/Card";
import {
  AdminApi,
  CertificateTemplateRef,
  LinkedCertificateTemplateUsage,
  NamespaceKind,
} from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";
import { AgentConfigurationForm } from "./AgentConfigurationForm";
import {
  ApplicationServicePrincipalLink,
  DeviceServicePrincipalLink,
} from "./DeviceServicePrincipalLink";
import { NamespaceContext } from "./NamespaceContext";
import { RefTableColumn, RefsTable } from "./RefsTable";
import { InputField } from "./InputField";
import { Button } from "../components/Button";
import Select, { SelectItem } from "../components/Select";

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
          <Select
            items={selectItems}
            label="Usage"
            selected={selectedUsage}
            setSelected={setSelectedUsage}
          />
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

export default function NamespacePage() {
  const { namespaceId } = useParams() as {
    profileType: NamespaceKind;
    namespaceId: string;
  };

  const { nsInfo } = React.useContext(NamespaceContext);
  const nsType = nsInfo?.type;
  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{nsType}</div>
      {(nsType === NamespaceKind.NamespaceKindCaRoot ||
        nsType === NamespaceKind.NamespaceKindCaInt ||
        nsType == NamespaceKind.NamespaceKindServicePrincipal ||
        nsType == NamespaceKind.NamespaceKindDevice) && (
        <Card>
          <CardSection>
            <CertificateTemplatesList
              nsType={nsType}
              namespaceId={namespaceId}
            />
          </CardSection>
        </Card>
      )}
      {nsType === NamespaceKind.NamespaceKindDevice && (
        <DeviceServicePrincipalLink namespaceId={namespaceId} />
      )}
      {nsType === NamespaceKind.NamespaceKindApplication && (
        <ApplicationServicePrincipalLink namespaceId={namespaceId} />
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
