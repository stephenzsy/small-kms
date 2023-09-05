import { Link } from "react-router-dom";
import { WellknownId } from "../constants";
import { nsDisplayNames } from "./displayConstants";

interface NamespaceProp {
  title: string;
  ids: string[];
}

const namespaces = {
  rootCa: {
    title: "Root Certificate Authorities",
    ids: [WellknownId.nsRootCa, WellknownId.nsTestRootCa],
  },
  intCa: {
    title: "Intermediate Certificate Authorities",
    ids: [WellknownId.nsIntCaIntranet, WellknownId.nsTestIntCa],
  },
};

function PolicySection(props: { namespace: NamespaceProp }) {
  const { namespace } = props;
  const { title, ids } = namespace;
  return (
    <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
      <h2 className="text-lg font-semibold">{title}</h2>
      <div className="mt-8 flow-root">
        <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
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
                  <th scope="col" className="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    <span className="sr-only">Edit</span>
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {ids.map((id) => (
                  <tr key={id}>
                    <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                      {id}
                    </td>
                    <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                      {nsDisplayNames[id] ?? ""}
                    </td>
                    <td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0 space-x-4">
                      <Link
                        to={`/admin/${id}/policies`}
                        className="text-indigo-600 hover:text-indigo-900"
                      >
                        Policies<span className="sr-only">, {id}</span>
                      </Link>
                      <a
                        href="#"
                        className="text-indigo-600 hover:text-indigo-900"
                      >
                        Certificates<span className="sr-only">, {id}</span>
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function AdminPage() {
  return (
    <>
      <PolicySection namespace={namespaces.rootCa} />
      <PolicySection namespace={namespaces.intCa} />
    </>
  );
}
