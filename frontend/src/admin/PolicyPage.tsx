import { XCircleIcon } from "@heroicons/react/24/solid";
import { useRequest, useSetState } from "ahooks";
import classNames from "classnames";
import { useId, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import {
  CertificateEnrollmentParameters,
  PolicyApi,
  WellKnownNamespaceId,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  AdminBreadcrumb,
  type BreadcrumbPageMetadata,
} from "./AdminBreadcrumb";
import { caBreadcrumPages } from "./CaPage";
import { v4 as uuidv4 } from "uuid";

const titleDisplayNames: Partial<Record<WellKnownNamespaceId, string>> = {
  [WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA]:
    "Cert enrollment policy: Root CA",
  [WellKnownNamespaceId.WellKnownNamespaceIDStr_IntCAService]:
    "Cert enrollment policy: Intermediate CA - Services",
  [WellKnownNamespaceId.WellKnownNamespaceIDStr_IntCAIntranet]:
    "Cert enrollment policy: intermediate CA - Intranet",
};

const sampleParams: CertificateEnrollmentParameters = {
  issuerId: uuidv4(),
  keyParameters: {
    kty: "RSA",
    size: 2048,
  },
  validity: "720h",
};

export default function PolicyPage() {
  const titleId = useId();
  const client = useAuthedClient(PolicyApi);
  const { namespaceId } = useParams();

  const {
    loading: createCertInProgress,
    data: policy,
    run: requestPolicy,
  } = useRequest(
    async (p?: CertificateEnrollmentParameters) => {
      if (p) {
        return client.putPolicyCertEnrollV1({
          namespaceId: namespaceId!,
          certificateEnrollmentParameters: p,
        });
      } else {
        return client.getPolicyCertEnrollV1({
          namespaceId: namespaceId!,
        });
      }
    },
    {
      refreshDeps: [namespaceId],
    }
  );

  const [paramsJsonStr, setParamsJsonStr] = useState("");

  const breadcrumPages: BreadcrumbPageMetadata[] = useMemo(() => {
    switch (namespaceId) {
      case WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA:
      case WellKnownNamespaceId.WellKnownNamespaceIDStr_IntCAService:
      case WellKnownNamespaceId.WellKnownNamespaceIDStr_IntCAIntranet:
        return [
          ...caBreadcrumPages,
          { name: "Create certificate authority", to: "#" },
        ];
    }
    return [{ name: "Create certificate", to: "#" }];
  }, [namespaceId]);

  return (
    <>
      <AdminBreadcrumb pages={breadcrumPages} />
      <form
        className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow"
        onSubmit={(e) => {
          e.preventDefault();
          requestPolicy(JSON.parse(paramsJsonStr));
        }}
      >
        <h1 id={titleId} className="px-4 py-5 sm:px-6 text-lg font-semibold">
          {titleDisplayNames[namespaceId as WellKnownNamespaceId] ??
            "Certificate enrollment policy"}
        </h1>

        <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
          <div className="border-neutral-900/10 space-y-6">
            <h2 className="text-base font-semibold leading-7 text-gray-900">
              Current policy document
            </h2>
            <div className="mt-2">
              {policy ? (
                <pre>{JSON.stringify(policy, undefined, 2)}</pre>
              ) : (
                <div>
                  No policy found. Sample:{" "}
                  <pre>{JSON.stringify(sampleParams, undefined, 2)}</pre>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
          <div className="border-neutral-900/10 space-y-6">
            <label
              htmlFor="policy-doc"
              className="text-base font-semibold leading-7 text-gray-900"
            >
              Policy document
            </label>
            <div className="mt-2">
              <textarea
                rows={10}
                id="policy-doc"
                className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                value={paramsJsonStr}
                onChange={(e) => setParamsJsonStr(e.target.value)}
              />
            </div>
          </div>
        </div>

        <div className="bg-neutral-50 px-4 py-4 sm:px-6 flex align-center justify-end gap-x-6">
          <button
            type="button"
            className="rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
            onClick={() => requestPolicy()}
          >
            Reset
          </button>
          <button
            type="submit"
            disabled={createCertInProgress}
            className={classNames(
              "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
              "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
            )}
          >
            {createCertInProgress ? "Creating..." : "Create"}
          </button>
        </div>
      </form>
    </>
  );
}
