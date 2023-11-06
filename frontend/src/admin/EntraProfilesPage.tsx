import { useRequest } from "ahooks";
import { Card, Table, TableColumnType, Typography } from "antd";
import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { Link } from "../components/Link";
import {
  AdminApi,
  NamespaceKind,
  ProfileRef,
  ResourceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { ImportProfileForm } from "./ImportProfileForm";

function useColumns(nsKind: NamespaceKind) {
  return useMemo<TableColumnType<ProfileRef>[]>(
    () => [
      {
        title: "Name",
        render: (r: ProfileRef) => <span className="font-mono">{r.id}</span>,
      },
      {
        title: "Display name",
        render: (r: ProfileRef) => r.displayName,
      },
      {
        title: "Actions",
        render: (r: ProfileRef) => <Link to={`/entra/${nsKind}/${r.id}`}>View</Link>,
      },
    ],
    [nsKind]
  );
}

export default function EntraProfilesPage() {
  const { nsKind } = useParams<{ nsKind: NamespaceKind }>();
  const adminApi = useAuthedClient(AdminApi);

  const { data: profiles, run: listProfiles } = useRequest(
    () => {
      return adminApi.listProfiles({
        profileResourceKind: nsKind as ResourceKind,
      });
    },
    {
      refreshDeps: [nsKind],
      ready:
        nsKind === NamespaceKind.NamespaceKindGroup ||
        nsKind === NamespaceKind.NamespaceKindUser ||
        nsKind === NamespaceKind.NamespaceKindServicePrincipal,
    }
  );

  const rootColumns = useColumns(NamespaceKind.NamespaceKindGroup);
  //   const onProfileUpsert: ProfileTypeMapRecord<() => void> = useMemo(() => {
  //     return {
  //       [ResourceKind.ProfileResourceKindRootCA]: listRootCAs,
  //       [ResourceKind.ProfileResourceKindIntermediateCA]: listIntermediateCAs,
  //     };
  //   }, [listRootCAs]);

  return (
    <>
      <Typography.Title>{nsKind}</Typography.Title>
      <Card title="Manage profiles">
        <Table<ProfileRef>
          columns={rootColumns}
          dataSource={profiles}
          rowKey={(r) => r.id}
        />
      </Card>
      {/* <Card title="Intermediate certificate authorities">
          <Table<ProfileRef>
            columns={intColumns}
            dataSource={intermediateCAs}
            rowKey={(r) => r.id}
          />
        </Card> */}

      {nsKind && (
        <Card title="Create certificate authority profile">
          <ImportProfileForm
            profileKind={nsKind as ResourceKind}
            onCreated={listProfiles}
          />
        </Card>
      )}
    </>
  );
}
