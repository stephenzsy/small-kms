import { Switch } from "@headlessui/react";
import { useRequest } from "ahooks";
import classNames from "classnames";
import { useId, useMemo, useState } from "react";
import { Link, generatePath, useMatch, useMatches } from "react-router-dom";
import { CertsApi, TestNamespaceId, WellKnownNamespaceId } from "../generated";
import { useCertsApi } from "../utils/useCertsApi";
import { AdminBreadcrumb, BreadcrumbPageMetadata } from "./AdminBreadcrumb";
import { useCaList } from "./useCaList";
import { RouteIds } from "../route-constants";

export const caBreadcrumPages: BreadcrumbPageMetadata[] = [
  { name: "CA", to: "/admin/ca" },
];
export const testcaBreadcrumPages: BreadcrumbPageMetadata[] = [
  { name: "Test CA", to: "/admin/testca" },
];

function CaSection({
  manageEnabled,
  certsApi,
  namespaceId,
  title,
  createButtonLabel,
}: {
  manageEnabled: boolean;
  certsApi: CertsApi;
  namespaceId: string;
  title: React.ReactNode;
  createButtonLabel?: string;
}) {
  const titleId = useId();
  const { data: caCerts } = useCaList(certsApi, namespaceId);
  return (
    <section className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow">
      <h2 id={titleId} className="px-4 py-5 sm:px-6 text-lg font-semibold">
        {title}
      </h2>
      {caCerts &&
        (caCerts.length > 0 ? (
          <ol className="divide-y divide-neutral-200" aria-labelledby={titleId}>
            {caCerts.map((cert) => (
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
                      namespaceId,
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
            No certificate found
          </div>
        ))}
      {manageEnabled && (
        <div className="bg-neutral-50 px-4 py-4 sm:px-6 flex align-center">
          <Link
            to={generatePath("/admin/cert/:namespaceId/new", {
              namespaceId,
            })}
            className="rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            {createButtonLabel || "Create certificate"}
          </Link>
        </div>
      )}
    </section>
  );
}

export default function AdminCaPage() {
  const client = useCertsApi();
  const [manageEnabled, setEnabled] = useState(false);
  const matches = useMatches();
  const isTest = useMemo(
    () => matches.some((m) => m.id === RouteIds.adminTestCa),
    [matches]
  );

  return (
    <>
      <AdminBreadcrumb
        pages={isTest ? testcaBreadcrumPages : caBreadcrumPages}
      />
      <section className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow px-4 py-5 sm:px-6">
        <Switch.Group as="div" className="flex items-center justify-between">
          <span className="flex flex-grow flex-col">
            <Switch.Label
              as="span"
              className="text-sm font-medium leading-6 text-gray-900"
              passive
            >
              Enable management
            </Switch.Label>
          </span>
          <Switch
            checked={manageEnabled}
            onChange={setEnabled}
            className={classNames(
              manageEnabled ? "bg-indigo-600" : "bg-gray-200",
              "relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2"
            )}
          >
            <span className="sr-only">Enable management</span>
            <span
              aria-hidden="true"
              className={classNames(
                manageEnabled ? "translate-x-5" : "translate-x-0",
                "pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
              )}
            />
          </Switch>
        </Switch.Group>
      </section>
      <CaSection
        certsApi={client}
        manageEnabled={manageEnabled}
        namespaceId={
          isTest
            ? TestNamespaceId.TestNamespaceIDStr_RootCA
            : WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA
        }
        title="Root CA"
        createButtonLabel="Create root CA"
      />
      <CaSection
        certsApi={client}
        manageEnabled={manageEnabled}
        namespaceId={WellKnownNamespaceId.WellKnownNamespaceIDStr_IntCAService}
        title="Intermediate CA - Services"
        createButtonLabel="Create intermediate CA"
      />
      <CaSection
        certsApi={client}
        manageEnabled={manageEnabled}
        namespaceId={WellKnownNamespaceId.WellKnownNamespaceIDStr_IntCAIntranet}
        title="Intermediate CA - Intranet"
        createButtonLabel="Create intermediate CA"
      />
    </>
  );
}
