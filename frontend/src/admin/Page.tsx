import { useRequest } from "ahooks";
import { Link } from "react-router-dom";
import { WellknownId } from "../constants";
import { AdminApi, DirectoryApi, NamespaceTypeShortName } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { RefsTable } from "./RefsTable";

const namespaceIds = {
  rootCa: [WellknownId.nsRootCa, WellknownId.nsTestRootCa],
  intCa: [WellknownId.nsIntCaIntranet, WellknownId.nsTestIntCa],
};

export default function AdminPage() {
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
        [NamespaceTypeShortName.NSType_ServicePrincipal]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_ServicePrincipal,
          }),
        [NamespaceTypeShortName.NSType_Group]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_Group,
          }),
        [NamespaceTypeShortName.NSType_Device]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_Device,
          }),
        [NamespaceTypeShortName.NSType_User]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_User,
          }),
        [NamespaceTypeShortName.NSType_Application]:
          await adminApi.listNamespacesByTypeV2({
            namespaceType: NamespaceTypeShortName.NSType_Application,
          }),
      };
    },
    {
      refreshDeps: [],
    }
  );

  return (
    <>
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_RootCA]}
        title="Root Certificate Authorities"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_RootCA}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_IntCA]}
        title="Intermediate Certificate Authorities"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_IntCA}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_ServicePrincipal]}
        title="Service Principals"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_ServicePrincipal}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_Group]}
        title="Groups"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_Group}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_Device]}
        title="Devices"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_Device}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_User]}
        title="Users"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_User}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
      <RefsTable
        items={allNs?.[NamespaceTypeShortName.NSType_Application]}
        title="Applications"
        refActions={(ref) => (
          <Link
            to={`/admin/${NamespaceTypeShortName.NSType_Application}/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        itemTitleMetadataKey="displayName"
      />
    </>
  );
}
