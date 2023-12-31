import { useContext } from "react";
import { NamespaceKind } from "../generated";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { CertificatesTableCard } from "./tables/CertificatesTableCard";

// function useGroupMemberOfColumns() {
//   return useMemo((): ColumnsType<ResourceReference> => {
//     return [
//       {
//         title: "ID",
//         render: (r: ResourceReference) => (
//           <span className="font-mono">{r.id}</span>
//         ),
//       },
//       {
//         title: "Actions",
//         render: (r: ResourceReference) => (
//           <Link to={`/entra/group/${r.id}`}>View</Link>
//         ),
//       },
//     ];
//   }, []);
// }

export default function NamespacePage() {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  // const adminApi = useAuthedClient(AdminApi);
  // const { data: groupMemberOf } = useRequest(
  //   async () => {
  //     return await adminApi.listGroupMemberOf({
  //       namespaceId,
  //       namespaceKind,
  //     });
  //   },
  //   {
  //     refreshDeps: [namespaceId, namespaceKind],
  //     ready: namespaceKind === NamespaceKind.NamespaceKindUser,
  //   }
  // );

  // const groupMemberOfColumns = useGroupMemberOfColumns();

  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{namespaceKind}</div>
      {namespaceKind === NamespaceKind.NamespaceKindUser && (
        <>
          <CertificatesTableCard />
          {/* <Card title="Listed group memberships">
            <Table<ResourceReference>
              dataSource={groupMemberOf}
              columns={groupMemberOfColumns}
              rowKey="id"
            />
          </Card> */}
        </>
      )}
    </>
  );
}
