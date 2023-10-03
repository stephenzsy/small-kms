import { Link, useParams } from "react-router-dom";
import { WellknownId, uuidNil } from "../constants";
import { nsDisplayNames } from "./displayConstants";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  AdminApi,
  DirectoryApi,
  NamespaceRef,
  NamespaceType,
  NamespaceTypeShortName,
  Ref as RRef,
} from "../generated";
import { useRequest } from "ahooks";
import { RefTableColumn, RefsTable, displayNameColumn } from "./RefsTable";
import { useMemo } from "react";
import { DeviceServicePrincipalLink } from "./DeviceServicePrincipalLink";
import { v5 as uuidv5 } from "uuid";
import { DeviceGroupInstall } from "./DeviceGroupInstall";
interface CreateDefaultLinkItem {
  id: string;
  title: string;
}

const groupSpDefaultTemplateName =
  "default-service-principal-client-credential";

function getGroupDefaultTemplateId(name: string) {
  return uuidv5(
    `https://example.com/#microsoft.graph.group/certificate-templates/${name}`,
    uuidv5.URL
  );
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
        namespaceType: nsType,
        namespaceId,
      });
    },
    {
      refreshDeps: [nsType, namespaceId],
    }
  );

  const createDefaultLinkItems = useMemo<CreateDefaultLinkItem[]>(() => {
    switch (nsType) {
      case NamespaceTypeShortName.NSType_RootCA:
      case NamespaceTypeShortName.NSType_IntCA:
      case NamespaceTypeShortName.NSType_ServicePrincipal:
        return [
          { id: uuidNil, title: "Create/update default certificate template" },
        ];
      case NamespaceTypeShortName.NSType_Group:
        const spTemplateId = getGroupDefaultTemplateId(
          groupSpDefaultTemplateName
        );
        return [
          {
            id: spTemplateId,
            title:
              "Create/update default credential template - service principal AAD client credential",
          },
        ];
    }
    return [];
  }, [data, nsType, namespaceId]);

  return (
    <RefsTable
      items={data}
      title="Certificate Templates"
      tableActions={
        <div className="flex gap-4">
          {createDefaultLinkItems.map((item, i) => (
            <Link
              key={item.id}
              to={`/admin/${nsType}/${namespaceId}/certificate-templates/${item.id}`}
              className="text-indigo-600 hover:text-indigo-900"
            >
              {item.title}
            </Link>
          ))}
        </div>
      }
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
          to={`/admin/${ref.namespaceType}/${ref.namespaceId}/certificate-templates/${ref.id}`}
          className="text-indigo-600 hover:text-indigo-900"
        >
          View
        </Link>
      )}
    />
  );
}

export default function NamespacePage() {
  const { nsType, namespaceId } = useParams() as {
    nsType: NamespaceTypeShortName;
    namespaceId: string;
  };
  const client = useAuthedClient(DirectoryApi);
  const adminApi = useAuthedClient(AdminApi);
  const { data: allNs } = useRequest(
    async () => {
      return {
        [NamespaceTypeShortName.NSType_RootCA]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_RootCA,
          }),
        [NamespaceTypeShortName.NSType_IntCA]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_IntCA,
          }),
      };
    },
    {
      refreshDeps: [],
    }
  );
  return (
    <>
      <h1>
        {nsType}/{namespaceId}
      </h1>
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
