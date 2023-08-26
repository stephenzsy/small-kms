import {
  AdminApi,
  CertificateRef,
  Configuration,
  ResponseError,
} from "@/generated";
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
    new Configuration({
      basePath: process.env.BACKEND_URL_BASE,
      accessToken: auth.accessToken,
    })
  );
  let certs: Array<CertificateRef> | undefined = undefined;
  try {
    certs = await api.listCertificates(
      {
        xCallerPrincipalName: auth.principalName!,
        xCallerPrincipalId: auth.principalId!,
        category: "root-ca",
      },
      { cache: "no-cache" }
    );
  } catch (e) {
    return <pre>{(e as ResponseError).message}</pre>;
  }

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
      <div className="overflow-hidden rounded-md bg-white shadow mt-6 border-gray-900 border">
        <div className="border-b border-gray-200 bg-white px-4 py-5 sm:px-6">
          <h3 className="text-base font-semibold leading-6 text-gray-900">
            Root CAs
          </h3>
        </div>
        <ul role="list" className="divide-y divide-gray-200">
          {certs.map((cert) => (
            <li key={cert.id} className="px-6 py-4">
              <dl>
                <div className="flex gap-x-2">
                  <dt className="font-medium">Common name</dt>
                  <dd>{cert.commonName}</dd>
                </div>
                <div className="flex gap-x-2">
                  <dt className="font-medium">Expires</dt>
                  <dd>{cert.notAfter.toLocaleString()}</dd>
                </div>
              </dl>
              <div className="mt-4 text-sm">
                These files should have been deployed through MDM
              </div>
              <div className="flex gap-4 mt-4">
                <a
                  className="rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                  href={`/api/admin/certificate/${cert.id}/download?format=pem`}
                  download={`root-ca-${cert.serialNumber}.pem`}
                >
                  Download .PEM file
                </a>
                <a
                  className="rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                  href={`/api/admin/certificate/${cert.id}/download?format=der`}
                  download={`root-ca-${cert.serialNumber}.der`}
                >
                  Download .DER file
                </a>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </main>
  );
}
