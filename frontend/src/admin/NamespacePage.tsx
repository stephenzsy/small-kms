import { useRequest } from "ahooks";
import React from "react";
import { Link, useParams } from "react-router-dom";
import { AdminApi, NamespaceTypeShortName } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { DeviceServicePrincipalLink } from "./DeviceServicePrincipalLink";
import { NamespaceContext } from "./NamespaceContext";
import { RefTableColumn, RefsTable, displayNameColumn } from "./RefsTable";

interface CreateDefaultLinkItem {
  id: string;
  title: string;
}

function CertificateTemplatesList({
  nsType,
  namespaceId,
}: {
  nsType: NamespaceTypeShortName;
  namespaceId: string;
}) {
  const adminApi = useAuthedClient(AdminApi);
  const { data } = useRequest(
    async () => {
      return await adminApi.listCertificateTemplatesV2({
        includeDefaultForType: nsType,
        namespaceId,
      });
    },
    {
      refreshDeps: [nsType, namespaceId],
    }
  );

  return (
    <RefsTable
      items={data}
      title="Certificate Templates"
      columns={
        [
          displayNameColumn,
          {
            header: "Active",
            metadataKey: "isActive",
            render: (isActive) => {
              return isActive ? "Yes" : "No";
            },
          } as RefTableColumn<"isActive">,
        ] as RefTableColumn[]
      }
      refActions={(ref) => (
        <Link
          to={`/admin/${ref.namespaceId}/certificate-templates/${ref.id}`}
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
    profileType: string;
    namespaceId: string;
  };

  const { nsInfo } = React.useContext(NamespaceContext);
  const nsType = nsInfo?.type;
  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{nsType}</div>
      {(nsType === "root-ca" ||
        nsType === NamespaceTypeShortName.NSType_IntCA ||
        nsType == NamespaceTypeShortName.NSType_ServicePrincipal ||
        nsType == NamespaceTypeShortName.NSType_Group) && (
        <CertificateTemplatesList nsType={nsType} namespaceId={namespaceId} />
      )}
      {nsType === NamespaceTypeShortName.NSType_Device && (
        <DeviceServicePrincipalLink namespaceId={namespaceId} />
      )}
    </>
  );
}
