import { Link } from "react-router-dom";
import { WellknownId } from "../constants";
import { nsDisplayNames } from "./displayConstants";
import { useAuthedClient } from "../utils/useCertsApi";
import { DirectoryApi, NamespaceRef, NamespaceType } from "../generated";
import { useRequest } from "ahooks";

const namespaceIds = {
  rootCa: [WellknownId.nsRootCa, WellknownId.nsTestRootCa],
  intCa: [WellknownId.nsIntCaIntranet, WellknownId.nsTestIntCa],
};

function PolicySection(props: {
  namespaces: Pick<NamespaceRef, "id" | "displayName">[] | undefined;
  title: string;
  showAdd?: boolean;
}) {
  const { namespaces, title, showAdd = false } = props;

  return (
    <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
      <div className="flex flex-row items-center justify-between">
        <h2 className="text-lg font-semibold">{title}</h2>
        {showAdd && (
          <Link
            to={`/admin/register`}
            className="text-indigo-600 font-semibold hover:text-indigo-900"
          >
            Add<span className="sr-only">namespace</span>
          </Link>
        )}
      </div>
      <div className="mt-8 flow-root">
        <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            {namespaces ? (
              <table className="min-w-full divide-y divide-gray-300">
                <thead>
                  <tr>
                    <th
                      scope="col"
                      className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-0"
                    >
                      ID
                    </th>
                    <th
                      scope="col"
                      className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                    >
                      Name
                    </th>
                    <th
                      scope="col"
                      className="relative py-3.5 pl-3 pr-4 sm:pr-0"
                    >
                      <span className="sr-only">Edit</span>
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {namespaces.map((ns) => (
                    <tr key={ns.id}>
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                        {ns.id}
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                        {ns.displayName}
                      </td>
                      <td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0 space-x-4">
                        <Link
                          to={`/admin/${ns.id}/policies`}
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          Policies<span className="sr-only">, {ns.id}</span>
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <div>Loading ...</div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}

export default function AdminPage() {
  const client = useAuthedClient(DirectoryApi);
  const { data: spNamespaces } = useRequest(
    () => {
      return client.listNamespacesV1({
        namespaceType: NamespaceType.NamespaceType_MsGraphServicePrincipal,
      });
    },
    { refreshDeps: [] }
  );

  const { data: gNamespaces } = useRequest(
    () => {
      return client.listNamespacesV1({
        namespaceType: NamespaceType.NamespaceType_MsGraphGroup,
      });
    },
    { refreshDeps: [] }
  );
  const { data: dNamespaces } = useRequest(
    () => {
      return client.listNamespacesV1({
        namespaceType: NamespaceType.NamespaceType_MsGraphDevice,
      });
    },
    { refreshDeps: [] }
  );
  const { data: uNamespaces } = useRequest(
    () => {
      return client.listNamespacesV1({
        namespaceType: NamespaceType.NamespaceType_MsGraphUser,
      });
    },
    { refreshDeps: [] }
  );
  const { data: aNamespaces } = useRequest(
    () => {
      return client.listNamespacesV1({
        namespaceType: NamespaceType.NamespaceType_MsGraphApplication,
      });
    },
    { refreshDeps: [] }
  );

  return (
    <>
      <PolicySection
        namespaces={namespaceIds.rootCa.map((id) => ({
          id,
          displayName: nsDisplayNames[id],
        }))}
        title="Root Certificate Authorities"
      />
      <PolicySection
        namespaces={namespaceIds.intCa.map((id) => ({
          id,
          displayName: nsDisplayNames[id],
        }))}
        title="Intermediate Certificate Authorities"
      />
      <PolicySection
        namespaces={spNamespaces}
        title="Service Principals"
        showAdd
      />
      <PolicySection namespaces={gNamespaces} title="Groups" showAdd />
      <PolicySection namespaces={dNamespaces} title="Devices" showAdd />
      <PolicySection namespaces={uNamespaces} title="Users" showAdd />
      <PolicySection namespaces={aNamespaces} title="Applications" showAdd />
    </>
  );
}
