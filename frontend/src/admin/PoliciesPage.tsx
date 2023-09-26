import { useRequest } from "ahooks";
import { useMemo } from "react";
import { Link, useParams } from "react-router-dom";
import { ErrorAlert } from "../components/ErrorAlert";
import { WellknownId } from "../constants";
import {
  DirectoryApi,
  NamespaceProfile,
  NamespaceType,
  PolicyApi,
} from "../generated";
import { PolicyRef } from "../generated/models/PolicyRef";
import { useAuthedClient } from "../utils/useCertsApi";
import { nsDisplayNames, policyTypeNames } from "./displayConstants";

function CreatePoliciesLinks({
  policyRefs,
  profile,
}: {
  policyRefs: PolicyRef[];
  profile: NamespaceProfile;
}) {
  const showCreateIntranet = useMemo(() => {
    if (profile.objectType !== "#microsoft.graph.servicePrincipal") {
      return false;
    }
    return !policyRefs.some((p) => p.id === WellknownId.nsIntCaIntranet);
  }, [policyRefs]);
  return (
    <>
      {showCreateIntranet && (
        <Link
          to={`/admin/${profile.id}/policies/${WellknownId.nsIntCaIntranet}`}
          className="text-indigo-600 hover:text-indigo-900"
        >
          Create Intranet Certificate Policy
          <span className="sr-only">, {WellknownId.nsIntCaIntranet}</span>
        </Link>
      )}
    </>
  );
}

export default function PoliciesPage() {
  const { namespaceId: _namespaceId } = useParams();
  const namespaceId = _namespaceId as string;
  const client = useAuthedClient(PolicyApi);
  const dirClient = useAuthedClient(DirectoryApi);

  const { data: policyRefs, error: fetchPoliciesError } = useRequest(
    () => {
      return client.listPoliciesV1({ namespaceId: namespaceId });
    },
    { refreshDeps: [] }
  );

  const { data: dirProfile } = useRequest(
    () => {
      return dirClient.getNamespaceProfileV1({ namespaceId: namespaceId });
    },
    { refreshDeps: [namespaceId] }
  );
  return (
    <>
      <h1 className="font-semibold text-4xl">Policies</h1>
      <div className="font-medium text-xl">
        {nsDisplayNames[namespaceId!] || dirProfile?.displayName || namespaceId}
      </div>
      {fetchPoliciesError && <ErrorAlert error={fetchPoliciesError} />}
      <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
        <div className="flow-root -mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            {policyRefs ? (
              policyRefs.length > 0 ? (
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
                        Type
                      </th>
                      <th
                        scope="col"
                        className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                      >
                        Status
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
                    {policyRefs.map((p) => (
                      <tr key={p.id}>
                        <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                          {p.id}
                        </td>
                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                          {policyTypeNames[p.policyType]}
                        </td>
                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                          {p.deleted ? "Disabled" : "Enabled"}
                        </td>
                        <td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0 space-x-4">
                          <Link
                            to={`/admin/${namespaceId}/policies/${p.id}`}
                            className="text-indigo-600 hover:text-indigo-900"
                          >
                            View
                          </Link>
                          <Link
                            to={`/admin/${namespaceId}/policies/${p.id}/latest-certificate`}
                            className="text-indigo-600 hover:text-indigo-900"
                          >
                            Certificate
                          </Link>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <div>No policy found</div>
              )
            ) : (
              <div>Loading ...</div>
            )}
          </div>
        </div>
      </section>
      {policyRefs &&
        dirProfile &&
        (dirProfile.objectType === NamespaceType.NamespaceType_BuiltInCaRoot ||
          dirProfile.objectType ===
            NamespaceType.NamespaceType_BuiltInCaInt) && (
          <>
            {policyRefs.some(
              (p) => p.id === WellknownId.defaultPolicyIdCertRequest
            ) ? null : (
              <Link
                to={`/admin/${namespaceId}/policies/${WellknownId.defaultPolicyIdCertRequest}`}
                className="text-indigo-600 hover:text-indigo-900"
              >
                Create Default Certificate Policy
              </Link>
            )}
          </>
        )}
      {policyRefs && dirProfile && (
        <div>
          <CreatePoliciesLinks policyRefs={policyRefs} profile={dirProfile} />
        </div>
      )}
    </>
  );
}
