import { useRequest } from "ahooks";
import { useId } from "react";
import { WellKnownNamespaceId } from "../generated";
import { useCertsApi } from "../utils/useCertsApi";
import { AdminBreadcrumb, BreadcrumbPageMetadata } from "./AdminBreadcrumb";
import { Link, generatePath } from "react-router-dom";

export const caBreadcrumPages: BreadcrumbPageMetadata[] = [
  { name: "CA", to: "/admin/ca" },
];

export default function AdminCaPage() {
  const titleId = useId();
  const client = useCertsApi();
  const { data: rootCaCerts } = useRequest(async () => {
    return client.listCertificatesV1({
      namespaceId: WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA,
    });
  }, {});
  return (
    <>
      <AdminBreadcrumb pages={caBreadcrumPages} />
      <section className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow">
        <h2 id={titleId} className="px-4 py-5 sm:px-6 text-lg font-semibold">
          Root Certificate Authority
        </h2>
        {rootCaCerts &&
          (rootCaCerts.length > 0 ? (
            <ol
              className="divide-y divide-neutral-200"
              aria-labelledby={titleId}
            >
              {rootCaCerts.map((cert) => (
                <li className="px-4 py-5 sm:p-6" key={cert.id}>
                  <dl>
                    <div className="flex gap-x-2">
                      <dt className="text-sm font-medium">Common Name</dt>
                      <dd>{cert.name}</dd>
                    </div>
                    <div className="flex gap-x-2">
                      <dt className="text-sm font-medium">Expiry</dt>
                      <dd>{cert.notAfter.toLocaleString()}</dd>
                    </div>
                  </dl>
                  <div className="mt-1">
                    <Link
                      to={generatePath("/admin/cert/:namespaceId/:certId", {
                        namespaceId:
                          WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA,
                        certId: cert.id,
                      })}
                      className="text-sm font-medium text-indigo-600 hover:text-indigo-500"
                    >
                      View details
                    </Link>
                  </div>
                </li>
              ))}
            </ol>
          ) : (
            <div className="px-4 py-5 sm:p-6 text-neutral-600">
              No root CA found
            </div>
          ))}
        <div className="bg-neutral-50 px-4 py-4 sm:px-6 flex align-center">
          <Link
            to={generatePath("/admin/cert/:namespaceId/new", {
              namespaceId: WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA,
            })}
            className="rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Create root CA
          </Link>
        </div>
      </section>
    </>
  );
}
