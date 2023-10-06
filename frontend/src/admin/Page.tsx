import { useRequest } from "ahooks";
import { Link } from "react-router-dom";
import { AdminApi, ProfileType } from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";
import {
  RefTableColumn,
  RefsTable,
  RefsTable3,
  displayNameColumn,
} from "./RefsTable";

export default function AdminPage() {
  const adminApi = useAuthedClient(AdminApi);
  const { data: allNs } = useRequest(
    async () => {
      return {
        [ProfileType.ProfileTypeRootCA]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeRootCA,
        }),
        [ProfileType.ProfileTypeIntermediateCA]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeIntermediateCA,
        }),
        [ProfileType.ProfileTypeServicePrincipal]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeServicePrincipal,
        }),
        [ProfileType.ProfileTypeGroup]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeGroup,
        }),
        [ProfileType.ProfileTypeDevice]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeDevice,
        }),
        [ProfileType.ProfileTypeUser]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeUser,
        }),
        [ProfileType.ProfileTypeApplication]: await adminApi.listProfiles({
          profileType: ProfileType.ProfileTypeApplication,
        }),
      };
    },
    {
      refreshDeps: [],
    }
  );

  return (
    <>
      {[
        ProfileType.ProfileTypeRootCA,
        ProfileType.ProfileTypeIntermediateCA,
      ].map((t) => (
        <RefsTable3
          key={t}
          items={allNs?.[t]}
          title="Root Certificate Authorities"
          refActions={(ref) => (
            <Link
              to={`/admin/${t}/${ref.id}`}
              className="text-indigo-600 hover:text-indigo-900"
            >
              View
            </Link>
          )}
          columns={[displayNameColumn] as RefTableColumn[]}
        />
      ))}
      <RefsTable3
        items={allNs?.[ProfileType.ProfileTypeServicePrincipal]}
        title="Service Principals"
        refActions={(ref) => (
          <Link
            to={`/admin/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        columns={[displayNameColumn] as RefTableColumn[]}
      />
      <RefsTable3
        items={allNs?.[ProfileType.ProfileTypeGroup]}
        title="Groups"
        refActions={(ref) => (
          <Link
            to={`/admin/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        columns={[displayNameColumn] as RefTableColumn[]}
      />
      <RefsTable3
        items={allNs?.[ProfileType.ProfileTypeDevice]}
        title="Devices"
        refActions={(ref) => (
          <Link
            to={`/admin/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        columns={[displayNameColumn] as RefTableColumn[]}
      />
      <RefsTable3
        items={allNs?.[ProfileType.ProfileTypeUser]}
        title="Users"
        refActions={(ref) => (
          <Link
            to={`/admin/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        columns={[displayNameColumn] as RefTableColumn[]}
      />
      <RefsTable3
        items={allNs?.[ProfileType.ProfileTypeApplication]}
        title="Applications"
        refActions={(ref) => (
          <Link
            to={`/admin/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
        columns={[displayNameColumn] as RefTableColumn[]}
      />
    </>
  );
}
