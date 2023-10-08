import { useRequest } from "ahooks";
import React from "react";
import { Link, useParams } from "react-router-dom";
import { AdminApi, CertificateTemplateRef, NamespaceKind } from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";
import { DeviceServicePrincipalLink } from "./DeviceServicePrincipalLink";
import { NamespaceContext } from "./NamespaceContext";
import {
  RefTableColumn,
  RefTableColumn3,
  RefsTable,
  RefsTable3,
  displayNameColumn,
} from "./RefsTable";

const subjectCnColumn: RefTableColumn3<CertificateTemplateRef> = {
  columnKey: "subjectCommonName",
  header: "Subject Common Name",
  render: (item) => item.subjectCommonName,
};
const enabledColumn: RefTableColumn3<CertificateTemplateRef> = {
  columnKey: "enabled",
  header: "Enabled",
  render: (item) => (item.metadata && !item.metadata.deleted ? "Yes" : "No"),
};

function CertificateTemplatesList({
  nsType,
  namespaceId,
}: {
  nsType: NamespaceKind;
  namespaceId: string;
}) {
  const adminApi = useAuthedClient(AdminApi);
  const { data } = useRequest(
    async () => {
      return await adminApi.listCertificateTemplates({
        profileId: namespaceId,
        profileType: nsType,
      });
    },
    {
      refreshDeps: [nsType, namespaceId],
    }
  );

  return (
    <RefsTable3
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
        <CertificateTemplatesList nsType={nsType} namespaceId={namespaceId} />
      )}
      {nsType === NamespaceKind.NamespaceKindDevice && (
        <DeviceServicePrincipalLink namespaceId={namespaceId} />
      )}
    </>
  );
}
