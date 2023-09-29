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
import { RefsTable } from "./RefsTable";
import { useMemo } from "react";
import { DeviceServicePrincipalLink } from "./DeviceServicePrincipalLink";

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

  const showCreateDefault = useMemo(() => {
    return !!data && !data.some((d) => d.id === uuidNil);
  }, [data]);

  return (
    <RefsTable
      items={data}
      title="Certificate Templates"
      tableActions={
        showCreateDefault && (
          <div>
            <Link
              to={`/admin/${nsType}/${namespaceId}/certificate-templates/${uuidNil}`}
              className="text-indigo-600 hover:text-indigo-900"
            >
              Create Default Certificate Template
            </Link>
          </div>
        )
      }
      itemTitleMetadataKey="displayName"
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
