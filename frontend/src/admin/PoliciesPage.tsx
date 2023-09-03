import { Link, useParams } from "react-router-dom";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  Policy,
  PolicyApi,
  PolicyType,
  ResponseError,
  TestNamespaceId,
  WellKnownNamespaceId,
} from "../generated";
import { useMemo } from "react";
import { useRequest } from "ahooks";
import {
  certRequestPolicyNames,
  isRootCANamespace,
  namespaceFriendlierNames,
} from "./displayConstants";

export default function PoliciesPage() {
  const { namespaceId } = useParams();
  const client = useAuthedClient(PolicyApi);
  const [fetchPolicyIds, catLabels] = useMemo<[string[], PolicyType[]]>(() => {
    switch (namespaceId) {
      case WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA:
      case TestNamespaceId.TestNamespaceIDStr_RootCA:
        return [[namespaceId], [PolicyType.PolicyType_CertRequest]];
    }
    return [[], []];
  }, [namespaceId]);
  const { data: fetchedPolicies, run: refresh } = useRequest(
    () => {
      return Promise.all(
        fetchPolicyIds.map(async (policyId): Promise<Policy | undefined> => {
          try {
            return await client.getPolicyV1({
              namespaceId: namespaceId!,
              policyId,
            });
          } catch (e) {
            if (e instanceof ResponseError && e.response.status === 404) {
              return undefined;
            }
            throw e;
          }
        })
      );
    },
    { refreshDeps: [fetchPolicyIds] }
  );
  return (
    <>
      <h1 className="font-semibold text-4xl">Policies</h1>
      <div className="font-medium text-xl">
        {namespaceFriendlierNames[namespaceId!] || namespaceId}
      </div>
      {catLabels.map((catLabel, i) => {
        return (
          <div
            key={i}
            className="divide-y space-y-4  divide-neutral-500 overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6"
          >
            <div>
              <h2 className="text-lg font-semibold mb-6">
                {certRequestPolicyNames[catLabel]}
              </h2>
              {!isRootCANamespace(namespaceId!) && (
                <dl>
                  <div>
                    <dt>CA Issuer Namespace</dt>
                    <dd>{namespaceFriendlierNames[catLabel]}</dd>
                  </div>
                </dl>
              )}
            </div>
            {fetchedPolicies && !fetchedPolicies[i] && (
              <div className="pt-4">Not found</div>
            )}
            <div className="pt-4">
              <Link
                to={`/admin/${namespaceId}/policies/${fetchPolicyIds[i]}`}
                className="text-indigo-600 hover:text-indigo-900 font-semibold"
              >
                Modify<span className="sr-only">, {fetchPolicyIds[i]}</span>
              </Link>
            </div>
          </div>
        );
      })}
    </>
  );
}
