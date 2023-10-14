import { useRequest } from "ahooks";
import { Link } from "react-router-dom";
import { AdminApi, NamespaceKind, ProfileRef } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { RefTableColumn, RefsTable } from "./RefsTable";
import { Card, CardSection } from "../components/Card";

const displayNameColumn: RefTableColumn<ProfileRef> = {
  columnKey: "displayName",
  header: "Display Name",
  render: (item) => item.displayName,
};

export default function AdminPage() {
  const adminApi = useAuthedClient(AdminApi);
  const { data: allNs } = useRequest(
    async () => {
      const results = {} as any;
      const l = [
        NamespaceKind.NamespaceKindCaRoot,
        NamespaceKind.NamespaceKindCaInt,
        NamespaceKind.NamespaceKindServicePrincipal,
        NamespaceKind.NamespaceKindGroup,
        NamespaceKind.NamespaceKindDevice,
        NamespaceKind.NamespaceKindUser,
        NamespaceKind.NamespaceKindApplication,
      ];
      for (const nsType of l) {
        results[nsType] = await adminApi.listProfiles({
          namespaceKind: nsType,
        });
      }
      return results;
    },
    {
      refreshDeps: [],
    }
  );

  return (
    <>
      {(
        [
          [NamespaceKind.NamespaceKindCaRoot, "Root Certificate Authorities"],
          [
            NamespaceKind.NamespaceKindCaInt,
            "Intermediate Certificate Authorities",
          ],
          [NamespaceKind.NamespaceKindServicePrincipal, "Service Principals"],
          [NamespaceKind.NamespaceKindGroup, "Groups"],
          [NamespaceKind.NamespaceKindDevice, "Devices"],
          [NamespaceKind.NamespaceKindUser, "Users"],
          [NamespaceKind.NamespaceKindApplication, "Applications"],
        ] as Array<[NamespaceKind, string]>
      ).map(([t, title]: [NamespaceKind, string]) => (
        <Card key={t}>
          <CardSection>
            <RefsTable
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
              tableActions={
                t == NamespaceKind.NamespaceKindApplication && (
                  <Link
                    to={`/admin/${t}/register`}
                    className="text-indigo-600 hover:text-indigo-900"
                  >
                    Register
                  </Link>
                )
              }
              columns={[displayNameColumn]}
            />
          </CardSection>
        </Card>
      ))}
    </>
  );
}
