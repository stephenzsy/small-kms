import { AdminApi, Configuration } from "@/generated";
import { getMsAuth } from "@/utils/aadAuthUtils";
import { ChevronRightIcon, HomeIcon } from "@heroicons/react/20/solid";
import { PlusIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

const pages = [
  { name: "Admin", href: "/admin", current: false },
  { name: "Certificate Authorities (CA)", href: "#", current: true },
];

export default async function CAIndex() {
  const auth = getMsAuth();
  const api = new AdminApi(
    new Configuration({ basePath: process.env.BACKEND_URL_BASE })
  );
  const certs = await api.listCACertificates(
    {
      xMsClientPrincipalName: auth.principalName!,
      xMsClientPrincipalId: auth.principalId!,
      xMsClientRoles: auth.isAdmin ? "App.Admin" : "",
    },
    { cache: "no-cache" }
  );
  return (
    <main>
      <nav className="flex" aria-label="Breadcrumb">
        <ol role="list" className="flex items-center space-x-4">
          <li>
            <div>
              <a href="/" className="text-gray-400 hover:text-gray-500">
                <HomeIcon
                  className="h-5 w-5 flex-shrink-0"
                  aria-hidden="true"
                />
                <span className="sr-only">Home</span>
              </a>
            </div>
          </li>
          {pages.map((page) => (
            <li key={page.name}>
              <div className="flex items-center">
                <ChevronRightIcon
                  className="h-5 w-5 flex-shrink-0 text-gray-400"
                  aria-hidden="true"
                />
                <Link
                  href={page.href}
                  className="ml-4 text-sm font-medium text-gray-500 hover:text-gray-700"
                  aria-current={page.current ? "page" : undefined}
                >
                  {page.name}
                </Link>
              </div>
            </li>
          ))}
        </ol>
      </nav>
      {certs.length === 0 && (
        <div className="text-center border-2 border-dashed border-gray-300 p-12 rounded-lg mt-6">
          <h3 className="mt-2 text-sm font-semibold text-gray-900">
            No CA certificate
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            Get started by creating a root certificate
          </p>
          <div className="mt-6">
            <Link
              className="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              href="/admin/ca/new/root"
            >
              <PlusIcon className="-ml-0.5 mr-1.5 h-5 w-5" aria-hidden="true" />
              Root CA
            </Link>
          </div>
        </div>
      )}
      <ul>
        {certs.map((cert) => (
          <li key={cert.id}>{JSON.stringify(cert, undefined, 2)}</li>
        ))}
      </ul>
    </main>
  );
}
