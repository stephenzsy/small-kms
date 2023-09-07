import { useRequest } from "ahooks";
import { useMemo } from "react";
import { Link, useParams } from "react-router-dom";
import { WellknownId } from "../constants";
import {
  DirectoryApi,
  Policy,
  PolicyApi,
  PolicyType,
  ResponseError,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  certRequestPolicyNames,
  isRootCaNamespace,
  nsDisplayNames,
} from "./displayConstants";
import { ErrorAlert } from "../components/ErrorAlert";
import { PolicyRef } from "../generated/models/PolicyRef";

const predefinedPolicyRefs: Record<string, PolicyRef[]> = {};

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
      {policyRefs &&
        (policyRefs.length == 0 ? (
          <div>No policy found</div>
        ) : (
          policyRefs?.map((policyRef) => {
            const policyId = policyRef.id;
            return (
              <div
                key={policyId}
                className="divide-y space-y-4  divide-neutral-500 overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6"
              >
                <div>
                  <h2 className="text-lg font-semibold mb-6">
                    {certRequestPolicyNames[policyRef.policyType]}
                  </h2>
                  {!isRootCaNamespace(namespaceId!) && (
                    <dl>
                      <div>
                        <dt>CA Issuer Namespace</dt>
                        <dd>{nsDisplayNames[policyId] ?? policyId}</dd>
                      </div>
                    </dl>
                  )}
                </div>
                <div className="pt-4">
                  <Link
                    to={`/admin/${namespaceId}/policies/${policyId}`}
                    className="text-indigo-600 hover:text-indigo-900 font-semibold"
                  >
                    Modify<span className="sr-only">, {policyId}</span>
                  </Link>
                </div>
              </div>
            );
          })
        ))}
    </>
  );
}
