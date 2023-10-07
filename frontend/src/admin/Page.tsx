import { useRequest } from "ahooks";
import { Link } from "react-router-dom";
import { AdminApi, ProfileRef, ProfileType } from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";
import {
  RefTableColumn,
  RefTableColumn3,
  RefsTable,
  RefsTable3,
} from "./RefsTable";

const displayNameColumn: RefTableColumn3<ProfileRef> = {
  columnKey: "displayName",
  header: "Display Name",
  render: (item) => item.displayName,
};

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
      {(
        [
          [ProfileType.ProfileTypeRootCA, "Root Certificate Authorities"],
          [
            ProfileType.ProfileTypeIntermediateCA,
            "Intermediate Certificate Authorities",
          ],
          [ProfileType.ProfileTypeServicePrincipal, "Service Principals"],
          [ProfileType.ProfileTypeGroup, "Groups"],
          [ProfileType.ProfileTypeDevice, "Devices"],
          [ProfileType.ProfileTypeUser, "Users"],
          [ProfileType.ProfileTypeApplication, "Applications"],
        ] as Array<[ProfileType, string]>
      ).map(([t, title]: [ProfileType, string]) => (
        <RefsTable3
          key={t}
          items={allNs?.[t]}
          title={title}
          refActions={(ref) => (
            <Link
              to={`/admin/${t}/${ref.id}`}
              className="text-indigo-600 hover:text-indigo-900"
            >
              View
            </Link>
          )}
          columns={[displayNameColumn]}
        />
      ))}
    </>
  );
}
